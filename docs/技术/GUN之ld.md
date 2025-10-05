本文将介绍 GNU ld 及其使用，参考[Using LD, the GNU linker](https://ftp.gnu.org/old-gnu/Manuals/ld-2.9.1/html_mono/ld.html)

## 1 Overview

ld将目标文件和静态库文件结合在一起，重定向他们的数据，合并符号引用。编译程序的最后一步通常是运行ld。

ld接受AT&T链接编辑命令语言语法（AT&T Link Editor Command Language syntax）格式的链接命令语言（Linker Command Language）文件，来显式地控制整个链接过程。

ld使用BFD库操作目标文件，这允许ld支持链接不同格式的目标文件。

许多链接器在发生错误后立即停止执行，而GNU ld尽可能继续执行，便于定位更多的错误。

## 2 命令行选项

链接器支持一大堆命令行选项，但是在任何特定的环境，都只有少量的选项被使用。例如，最基本的链接命令：

```bash
ld -o a.out /lib/crt0.o hello.o -lc
```

上述命令告诉ld，链接crt0.o，hello.o和libc.a或者libc.so产生一个可执行文件a.out。

### 2.1 选项顺序

一些选项可以出现在任何位置。然而涉及文件的命令，如-l或者-T，文件会在命令出现的位置被读入。保持和目标文件及其它文件选项的相对位置。后面可以看到，这些顺序是重要的。

### 2.2 重复选项

非文件选项重复，且其参数与之前的不同。或者没有影响，或者覆盖之前（命令行左侧）的参数。出现多次有意义的参数在后面会提到。

### 2.3 非选项参数

非选项参数是要被链接的目标文件或者静态库文件，他们可以出现任何位置。

### 2.4 链接脚本

链接脚本可以用非选项参数的方式设置，链接器会先将其当成目标文件格式，当不能识别目标文件的格式时，他将其当成一个链接脚本。这种方式的链接脚本，增大主链接器脚本，而不是替换。主链接器脚本是缺省链接脚本，或者-T指定的链接脚本。

### 2.5 选项参数格式

单字符选项，前面有且仅有一个破折号，参数或者选项连在一起，或者作为一个独立的参数，立即跟在选项后面。如"`-lc`" 或"`-l c`"。

多字符选项，选项名前有一个或两个破折号，如"`-trace-symbol`"与"`--trace-symbol`"是一样的。唯一的例外是，`o` 开头的多字符选项必须使用双破折号，如"`--omagic`"，这是为了避免和 `-o` 冲突。

多字符选项的参数同样有两种形式，或者使用=、或者作为一个独立的参数立即跟在选项后面，如"`--trace-symbol=foo`"与"`--trace-symbol foo`"是一样的。

### 2.6 gcc传递链接选项

如果链接器不是直接被调用，而是通过一个编译器驱动（如gcc），那么所有的链接器选项应当使用一个前缀"`-Wl`,"，或者其它编译器驱动支持的形式，如

```bash
gcc -Wl,--start-group foo.o bar.o -Wl,--end-group
```

这很重要，因为上层驱动可能默默地丢掉链接选项。

如果选项带参数，空格会导致上层驱动误解，这时可以使用不带空格形式，单字符选项连在一起，多字符选项使用=。如

```bash
gcc foo.o bar.o -Wl,-eENTRY -Wl,-Map=a.map
```

或者将整个链接选项用引号（我认为可以，手册没提到这种用法），如：

```bash
gcc foo.o bar.o "-Wl,-e ENTRY" "-Wl,-Map a.map"
```

## 3 部分参数

- `-e entry`, `--entry=entry`

使用entry符号，作为开始程序执行的入口，而不是缺省的入口点。如果没有entry符号，则链接器将其当做数字，作为入口地址。数字被翻译为10进制的；如果是0x开头，则被翻译为16进制的；如果是0开头，则是8进制。

- `-fini=name`

创建可执行文件和共享文件时，设置`DT_FINI`为`name`指定的地址。缺省，链接器使用`_fini`。

- `-h name`, `-soname=name`

创建共享目标文件时，指定`DT_SONAME`的值。

- `-init=name`

创建可执行文件或者共享文件时，设置`DT_INIT`为`name`指定的地址，缺省，链接器使用`_init`。

- `-l namespace`, `--library=namespec`

添加静态库或动态库，先搜索`libnamespec.so`的动态库，没有再搜索`libnamespec.a`的静态库。如果`namespec`是`:filename`的形式，则直接搜索叫`filename`文件名的库。链接器只会搜索库一次，也就是说，如果库定义了一个符号，而目标文件在方面，不会往前搜索。所以通常被依赖者要放在后面。可以指定同一个库多次。

- `-L searchdir`, `--library-path=searchdir`

指定库的搜索路径。如果使用`-T`指定了链接脚本，也指定了链接脚本的搜索路径。搜索路径是全局的，即使他们放在`-l`之后也同样生效。搜索的顺序，按照`-L`的顺序。

- `-o ouput`, `--output=output`

指定输出文件的名字，如果没有指定为`a.out`。

- `-q`, `--emit-relocs`

在完全链接的可执行文件中，保留可重定向section和内容。

- `-i`, `-r`, `--relocatable`

产生可重定向目标文件。

- `-s`, `--strip-all`

清掉输出文件的所有符号信息。

- `-T scriptfile`, `--script=scriptfile`

指定链接器脚本，如果不在当前目录，则尝试搜索-L指定的目录。

- `-x`, `--discard-all`

删除所有本地符号

- `-X`, `--discard-locals`

删除所有本地临时符号，这些符号通常以系统定义的本地符号前缀开始，如.L

- `-( archives -)`, `--start-group archives --end-group`

archives是共享文件的列表，它们可以是显式的名字或者-l选项。这些文件将被搜索多次来解析未定义符号，通常只有循环依赖的时候才需要这种方式。

- `-Bdynmic`, `-dy`, `-call_shared`

链接共享库，而不是静态库，影响后续的-l选项。

- `-Bstatic`, `-dn`, `-non_shared`, `-static`

链接静态库，而不是共享库，影响后续的-l选项。

- `-Ifile`, `--dynamic-linker=file`

指定动态链接器的名字。

- `--fatal-warnings`, `--no-fatal-warnings`

是/否将所有的警告当做错误。影响之后的选项。

- `-M`, `--print-map`, `-Map=file`

打印link map到stdout或者文件。

- `-nostdlib`

只搜索命令行或者脚本指定的链接路径。

- `-pie`, `--pic-executable`

创建位置独立可执行文件。

- `-rpath=dir`

指定运行时库搜索路径。

- `-shared`, `-Bshareable`

创建共享库。

- `--section-start=sectionname=org`, `-Tbss=org`, `-Tdata=org`, `-Ttext=org`

定位`section`到绝对地址。org必须是16进制，可以省略前导`0x`。

- `-Ttext-segment=org`, `-Trodata-segment=org`

设置`text`，`rodata`段的第一个字节到绝对地址。

## 4 crt*.o

> CRT: C Run-Time Libraries

- `crtbegin.o` ，`crtend.o`

包含全局构造、析构函数相关部分。

- `crt1.o`

包含程序的入口函数 `_start`，由它负责调用 `__libc_start_main` 初始化libc并且调用 `main` 函数进入真正的程序主体。

- `crti.o`， `crtn.o`

分别包含 `.init` 段和 `.finit` 段的一些辅助代码。

运行库在目标文件后引入 `.init`段（`main()`前执行）和 `.finit` 段（`main()`后执行）。

链接器将 `.init` 段与 `.finit` 段合并，并产生 `_init()` 与 `_finit()`。

![call_graph](../assets/callgraph.png)
