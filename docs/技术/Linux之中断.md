
简单梳理下Linux中的中断机制。

## 硬件中断和软件中断

中断分为硬件中断和软件中断

- 硬件中断：由硬件引发的中断，例如外部io设备引发的中断、CPU内部异常，区别于硬中断
- 软件中断：由软件引发的中断，例如int指令，区别于软中断

## 硬中断和软中断

中断执行过程通常分为上下半，以节省硬中断资源（也有简单的中断没有下半）

- 上半 硬中断，由硬件实现的中断机制，中断资源是有限的，需要尽快处理完，因此硬中断主要将软中断标志位置位，然后返回
- 下半 软中断，由软件实现的中断机制，不同的软中断标志由对应的守护程序进行轮询

### 硬中断

简单来说，硬中断的过程如下：

```
触发 --> CPU获得中断号 --> 通过硬布线解码得到处理程序地址 --> 跳转执行 --> 完成返回
※ 期间根据是否支持多级中断，还需设置中断使能位
```

可以看出，硬中断的资源和硬件资源相关，是十分宝贵的，但是其响应也相对快。

### 软中断

看看Linux 2.6.0代码：

```c
// 直接在init/main.c找softirq相关的初始化函数
static void do_pre_smp_initcalls(void)
{
	extern int spawn_ksoftirqd(void);
	node_nr_running_init();
	spawn_ksoftirqd(); //  <----- Here
}

// 注册了个回调函数
__init int spawn_ksoftirqd(void)
{
	cpu_callback(&cpu_nfb, CPU_ONLINE, (void *)(long)smp_processor_id()); // <-- Here
	register_cpu_notifier(&cpu_nfb);
	return 0;
}

// 当某个软中断触发，就通知
static int __devinit cpu_callback(struct notifier_block *nfb,
				  unsigned long action,
				  void *hcpu)
{
	int hotcpu = (unsigned long)hcpu;

	if (action == CPU_ONLINE) {
		if (kernel_thread(ksoftirqd, hcpu, CLONE_KERNEL) < 0) { // <---- Here
			printk("ksoftirqd for %i failed\n", hotcpu);
			return NOTIFY_BAD;
		}

		while (!per_cpu(ksoftirqd, hotcpu))
			yield();
 	}
	return NOTIFY_OK;
}

// 怎么触发？
// 主要是一个循环进行do_softirq()，再看看 do_softirq()
static int ksoftirqd(void * __bind_cpu)
{
	int cpu = (int) (long) __bind_cpu;

	daemonize("ksoftirqd/%d", cpu);
	set_user_nice(current, 19);
	current->flags |= PF_IOTHREAD;

	/* Migrate to the right CPU */
	set_cpus_allowed(current, cpumask_of_cpu(cpu));
	BUG_ON(smp_processor_id() != cpu);

	__set_current_state(TASK_INTERRUPTIBLE);
	mb();

	__get_cpu_var(ksoftirqd) = current;

	for (;;) {
		if (!local_softirq_pending())
			schedule();

		__set_current_state(TASK_RUNNING);

		while (local_softirq_pending()) {
			do_softirq();
			cond_resched();
		}

		__set_current_state(TASK_INTERRUPTIBLE);
	}
}

// 主要是对软中断向量表的各个位进行判断，并选择是否执行对应处理函数
asmlinkage void do_softirq(void)
{
	int max_restart = MAX_SOFTIRQ_RESTART;
	__u32 pending;
	unsigned long flags;

	if (in_interrupt())
		return;

	local_irq_save(flags);

	pending = local_softirq_pending();

	if (pending) {
		struct softirq_action *h;

		local_bh_disable();
restart:
		/* Reset the pending bitmask before enabling irqs */
		local_softirq_pending() = 0;

		local_irq_enable();

		h = softirq_vec; // 顾名思义：软中断向量表

		do {
			if (pending & 1) // 如果中断标志位置位
				h->action(h); // 调用处理函数
			h++;           // 下一个中断处理函数
			pending >>= 1; // 下一个中断标志位
		} while (pending);

		local_irq_disable();

		pending = local_softirq_pending();
		if (pending && --max_restart)
			goto restart;
		if (pending)
			wakeup_softirqd();
		__local_bh_enable();
	}

	local_irq_restore(flags);
}
// 看看软中断标志相关
typedef struct {
	unsigned int __softirq_pending;
	unsigned long idle_timestamp;
	unsigned int __nmi_count;	/* arch dependent */
	unsigned int apic_timer_irqs;	/* arch dependent */
} ____cacheline_aligned irq_cpustat_t;
extern irq_cpustat_t irq_stat[];		/* defined in asm/hardirq.h */
#define __IRQ_STAT(cpu, member)	((void)(cpu), irq_stat[0].member)
#define softirq_pending(cpu)	__IRQ_STAT((cpu), __softirq_pending)
#define local_softirq_pending()	softirq_pending(smp_processor_id())

// 看看软中断向量表相关结构
static struct softirq_action softirq_vec[32] __cacheline_aligned_in_smp;
struct softirq_action
{
	void	(*action)(struct softirq_action *);
	void	*data;
};
```

精简一下

```c
// 直接在init/main.c找softirq相关的初始化函数
static void do_pre_smp_initcalls(void)
{
	spawn_ksoftirqd(); //  <------ Here
}

// 注册了个回调函数
__init int spawn_ksoftirqd(void)
{
	cpu_callback(&cpu_nfb, CPU_ONLINE, (void *)(long)smp_processor_id()); // <--- Here
	register_cpu_notifier(&cpu_nfb);
	return 0;
}

// 当某个软中断触发，就通知
static int __devinit cpu_callback(struct notifier_block *nfb, unsigned long action, void *hcpu)
{
	kernel_thread(ksoftirqd, hcpu, CLONE_KERNEL); //   <------- Here
	return NOTIFY_OK;
}

// 怎么触发？
// 主要是一个循环进行do_softirq()，再看看 do_softirq()
static int ksoftirqd(void * __bind_cpu)
{
	for (;;) {
		while (local_softirq_pending()) {
			do_softirq();
			cond_resched();
		}
	}
}

// 主要是对软中断向量表的各个位进行判断，并选择是否执行对应处理函数
asmlinkage void do_softirq(void)
{
	__u32 pending = local_softirq_pending();
    struct softirq_action *h = softirq_vec; // 故名思意：软中断向量表
    do {
        if (pending & 1) // 如果中断标志位置位
            h->action(h); // 调用处理函数
        h++;           // 下一个中断处理函数
        pending >>= 1; // 下一个中断标志位
    } while (pending);
}

// 看看软中断标志相关
typedef struct {
	unsigned int __softirq_pending;
	unsigned long idle_timestamp;
	unsigned int __nmi_count;	/* arch dependent */
	unsigned int apic_timer_irqs;	/* arch dependent */
} ____cacheline_aligned irq_cpustat_t;
extern irq_cpustat_t irq_stat[];		/* defined in asm/hardirq.h */
#define local_softirq_pending()	irq_stat[0].__softirq_pending // 展开相关宏

// 看看软中断向量表相关结构
struct softirq_action
{
	void	(*action)(struct softirq_action *);
	void	*data;
};
static struct softirq_action softirq_vec[32] __cacheline_aligned_in_smp;

```

可以发现，软中断标志位及向量表最多同时支持32项。下面是已有的：

```c
enum
{
	HI_SOFTIRQ=0,
	TIMER_SOFTIRQ,
	NET_TX_SOFTIRQ,
	NET_RX_SOFTIRQ,
	SCSI_SOFTIRQ,
	TASKLET_SOFTIRQ
};
```

软中断系统提供下面的接口：

```c
asmlinkage void do_softirq(void);
// 注册中断
extern void open_softirq(int nr, void (*action)(struct softirq_action*), void *data);
extern void softirq_init(void);
// 触发中断：实际上就是将对应的中断标志位置位
#define __raise_softirq_irqoff(nr) do { local_softirq_pending() |= 1UL << (nr); } while (0)
extern void FASTCALL(raise_softirq_irqoff(unsigned int nr));
extern void FASTCALL(raise_softirq(unsigned int nr));

#ifndef invoke_softirq
#define invoke_softirq() do_softirq()
#endif
```
