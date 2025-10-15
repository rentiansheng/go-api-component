# PKG 公共包

包含项目的通用工具和配置组件。

## 组件结构

```
pkg/
├── config/          # 配置管理
│   ├── config.go    # 配置处理核心逻辑
│   └── README.md
└── logger/          # 日志工具
    ├── log.go       # 日志工具实现
    └── README.md
```

## 组件说明

### Config 配置管理

基于 `viper` 的配置管理系统，支持：
- 多种配置文件格式（JSON、YAML、TOML 等）
- 环境变量绑定
- 默认值设置
- 配置热重载

### Logger 日志工具

提供结构化日志功能，基于 `logrus` 实现：
- 多种日志级别
- JSON 格式输出
- 日志轮转
- 自定义字段

## 使用示例

### 配置管理使用

```go
import "github.com/rentiansheng/go-api-component/pkg/config"

// 创建配置处理器
cfg := config.NewConfigHandlerWithDefaults("app")

// 加载配置
if err := cfg.Load(); err != nil {
    log.Fatal("加载配置失败:", err)
}

// 读取配置
serverPort := cfg.GetString("server.port")
dbHost := cfg.GetString("database.host")
```

### 日志工具使用

```go
import "github.com/rentiansheng/go-api-component/pkg/logger"

// 创建日志实例
log := logger.NewLogger()

// 记录日志
log.WithFields(logger.Fields{
    "userID": 123,
    "action": "login",
}).Info("用户登录成功")
```

详细使用说明请参考各子组件的 README 文件。