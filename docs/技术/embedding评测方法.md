## mteb
### 介绍
MTEB（Multilingual Text Embedding Benchmark）是一个用于评估多语言文本嵌入的基准测试。它提供了一组标准化的任务和数据集，以便研究人员可以比较不同的多语言嵌入模型的性能。

### 评测方法
MTEB评测方法包括以下几个步骤：
1. 选择任务：MTEB提供多种任务，也能够自定义任务。
2. 选择模型：提供多种常见模型，也可以使用自定义模型。
3. 运行评测：在选定的任务和模型上运行评测。
4. 分析结果：对评测结果进行处理。

示例：
```python
import json
from mteb import MTEB
import mteb
from hunyuan import HunyuanModel
from ollama_model import OllamaModel

benchs = [
    {"model": HunyuanModel(), "output_file": "results_hunyuan.json"},
    {"model": OllamaModel("bge-m3"), "output_file": "results_bge_m3.json"},
    {"model": OllamaModel("nomic-embed-text"), "output_file": "results_nomic_embed_text.json"},
]

for bench in benchs:
    results = None
    # Select model
    model = bench["model"]

    # Select tasks
    evaluation = MTEB(tasks=["CodeFeedbackMT"])

    # evaluate
    results = evaluation.run(model)
    
    # results is NamedTuple, encode to json and save to file
    results_json = json.dumps({"results": [results[i].to_dict() for i in range(len(results))]}, indent=4)
    with open(bench["output_file"], "w") as f:
        f.write(results_json)
```

### 自定义模型
模型需要实现一个`encode`方法，该方法接受一个字符串列表作为输入，并返回一个二维浮点数列表作为输出。每个输入字符串对应一个嵌入向量。
下面是一个简单的自定义模型示例：
```python
class OllamaModel:
    def __init__(self, model_name: str = "bge-m3") -> None:
        self.model_name = model_name
    def encode(
        self,
        inputs: List[str],
        batch_size=50,
        **kwargs,
    ) -> np.ndarray:
        # split batches
        batches = [inputs[i:min(i + batch_size, len(inputs))] for i in range(0, len(inputs), batch_size)]
        embeddings = []
        for batch in batches:
            resp = ollama.embed(self.model_name, batch)
            for item in resp["embeddings"]:
                embeddings.append(item)
        return np.array(embeddings)
```
