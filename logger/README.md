# Logger 日志组件

提供日志记录功能的核心组件。

## 功能特性

- 支持多种日志级别
- 结构化日志输出
- 日志轮转和归档
- 自定义日志格式

## 使用方法

```go
import "github.com/rentiansheng/go-api-component/logger"

// 创建日志实例
log := logger.NewLogger()

// 记录不同级别的日志
log.Info("应用启动")
log.Warn("警告信息")
log.Error("错误信息")
log.Debug("调试信息")
```

## 配置

日志组件支持通过配置文件进行自定义：

```yaml
logger:
  level: info
  format: json
  output: stdout
  file:
    filename: app.log
    max_size: 100
    max_age: 30
    max_backups: 10
```