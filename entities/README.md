# Entities 目录

`entities` 按领域职责组织数据模型，目录之间保持单向、清晰的依赖关系：

```text
entities/
├── action/          # 调用后端、拓扑和验证服务时使用的 JSON 请求参数
├── configuration/   # 单次仿真的配置
├── event/           # 仿真调度事件
├── topology/        # 节点、链路和拓扑启动参数
├── sec_path_mab/    # SecPathMAB 专用拓扑参数
├── proto/           # protobuf 源文件及生成脚本
├── types/           # protobuf 生成的 Go 类型
└── entities.go      # 旧导入路径的兼容别名
```

新代码应直接导入具体子包。例如：

```go
import (
	"chain_simulation/entities/action"
	"chain_simulation/entities/event"
	"chain_simulation/entities/topology"
)
```

根包中的类型别名用于保持现有实验代码兼容，避免目录整理改变运行行为。新增实体时应放入对应子包，不再放到根目录。
