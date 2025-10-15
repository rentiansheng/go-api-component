# Logger 日志工具

基于 `logrus` 和 `lumberjack` 的结构化日志系统，提供完整的日志记录和管理功能。

## 特性

- 📝 结构化日志输出（JSON/Text 格式）
- 🗂️ 日志轮转和归档（基于 lumberjack）
- 🎯 多种日志级别支持
- 🏷️ 自定义字段支持
- 🔧 灵活的配置选项
- 🎨 Hook 机制扩展

## 使用方法

### 基础使用

```go
import "github.com/rentiansheng/go-api-component/pkg/logger"

// 创建日志实例
log := logger.NewLogger()

// 记录不同级别的日志
log.Debug("调试信息")
log.Info("普通信息")
log.Warn("警告信息") 
log.Error("错误信息")
log.Fatal("致命错误") // 会调用 os.Exit(1)
log.Panic("恐慌错误") // 会调用 panic()
```

### 带字段的日志

```go
// 使用 WithField 添加单个字段
log.WithField("userID", 12345).Info("用户登录成功")

// 使用 WithFields 添加多个字段
log.WithFields(logger.Fields{
    "userID":   12345,
    "username": "john_doe",
    "ip":       "192.168.1.100",
    "action":   "login",
}).Info("用户操作记录")

// 链式调用
log.WithField("module", "auth").
    WithField("function", "authenticate").
    WithError(err).
    Error("认证失败")
```

### 错误日志

```go
import "errors"

err := errors.New("数据库连接失败")

// 记录错误
log.WithError(err).Error("操作失败")

// 带上下文的错误日志
log.WithFields(logger.Fields{
    "database": "postgres",
    "host":     "localhost",
    "port":     5432,
}).WithError(err).Error("数据库连接失败")
```

## 配置选项

### 日志级别

```go
// 设置日志级别
log.SetLevel(logger.DebugLevel)   // 调试级别
log.SetLevel(logger.InfoLevel)    // 信息级别（默认）
log.SetLevel(logger.WarnLevel)    // 警告级别
log.SetLevel(logger.ErrorLevel)   // 错误级别
log.SetLevel(logger.FatalLevel)   // 致命错误级别
log.SetLevel(logger.PanicLevel)   // 恐慌级别
```

### 输出格式

```go
// JSON 格式输出（推荐用于生产环境）
log.SetFormatter(&logger.JSONFormatter{})

// 文本格式输出（推荐用于开发环境）
log.SetFormatter(&logger.TextFormatter{
    FullTimestamp: true,
    ForceColors:   true,
})
```

### 输出目标

```go
import "os"

// 输出到标准输出
log.SetOutput(os.Stdout)

// 输出到标准错误
log.SetOutput(os.Stderr)

// 输出到文件
file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
if err != nil {
    log.Fatal(err)
}
log.SetOutput(file)
```

## 日志轮转配置

使用 `lumberjack` 实现日志文件轮转：

```go
import "gopkg.in/natefinch/lumberjack.v2"

// 配置日志轮转
logRotate := &lumberjack.Logger{
    Filename:   "app.log",    // 日志文件路径
    MaxSize:    100,          // 单个文件最大 MB
    MaxAge:     30,           // 文件保留天数
    MaxBackups: 10,           // 最大备份文件数
    LocalTime:  true,         // 使用本地时间
    Compress:   true,         // 压缩备份文件
}

log.SetOutput(logRotate)
```

## 完整配置示例

```go
package main

import (
    "os"
    
    "github.com/rentiansheng/go-api-component/pkg/logger"
    "gopkg.in/natefinch/lumberjack.v2"
)

type LoggerConfig struct {
    Level      string `mapstructure:"level"`
    Format     string `mapstructure:"format"`
    Output     string `mapstructure:"output"`
    Filename   string `mapstructure:"filename"`
    MaxSize    int    `mapstructure:"max_size"`
    MaxAge     int    `mapstructure:"max_age"`
    MaxBackups int    `mapstructure:"max_backups"`
    Compress   bool   `mapstructure:"compress"`
}

func NewLogger(config LoggerConfig) *logger.Log {
    log := logger.NewLogger()
    
    // 设置日志级别
    switch config.Level {
    case "debug":
        log.SetLevel(logger.DebugLevel)
    case "info":
        log.SetLevel(logger.InfoLevel)
    case "warn":
        log.SetLevel(logger.WarnLevel)
    case "error":
        log.SetLevel(logger.ErrorLevel)
    default:
        log.SetLevel(logger.InfoLevel)
    }
    
    // 设置输出格式
    switch config.Format {
    case "json":
        log.SetFormatter(&logger.JSONFormatter{})
    case "text":
        log.SetFormatter(&logger.TextFormatter{
            FullTimestamp: true,
        })
    default:
        log.SetFormatter(&logger.JSONFormatter{})
    }
    
    // 设置输出目标
    switch config.Output {
    case "stdout":
        log.SetOutput(os.Stdout)
    case "stderr":
        log.SetOutput(os.Stderr)
    case "file":
        if config.Filename != "" {
            // 使用日志轮转
            logRotate := &lumberjack.Logger{
                Filename:   config.Filename,
                MaxSize:    config.MaxSize,
                MaxAge:     config.MaxAge,
                MaxBackups: config.MaxBackups,
                LocalTime:  true,
                Compress:   config.Compress,
            }
            log.SetOutput(logRotate)
        }
    default:
        log.SetOutput(os.Stdout)
    }
    
    return log
}

func main() {
    // 从配置创建日志实例
    config := LoggerConfig{
        Level:      "info",
        Format:     "json",
        Output:     "file",
        Filename:   "app.log",
        MaxSize:    100,
        MaxAge:     30,
        MaxBackups: 10,
        Compress:   true,
    }
    
    log := NewLogger(config)
    
    // 使用日志
    log.Info("应用启动")
    log.WithFields(logger.Fields{
        "version": "1.0.0",
        "env":     "production",
    }).Info("应用信息")
}
```

## 配置文件示例

### YAML 配置

```yaml
logger:
  level: info
  format: json
  output: file
  filename: logs/app.log
  max_size: 100
  max_age: 30
  max_backups: 10
  compress: true
```

### JSON 配置

```json
{
  "logger": {
    "level": "info",
    "format": "json", 
    "output": "file",
    "filename": "logs/app.log",
    "max_size": 100,
    "max_age": 30,
    "max_backups": 10,
    "compress": true
  }
}
```

## Hook 机制

添加自定义 Hook 来扩展日志功能：

```go
import "github.com/sirupsen/logrus"

// 自定义 Hook
type MyHook struct{}

func (hook *MyHook) Fire(entry *logrus.Entry) error {
    // 在这里可以添加自定义逻辑
    // 比如发送错误到监控系统、写入数据库等
    
    if entry.Level <= logrus.ErrorLevel {
        // 错误级别日志的特殊处理
        sendToMonitoring(entry)
    }
    
    return nil
}

func (hook *MyHook) Levels() []logrus.Level {
    return logrus.AllLevels
}

// 添加 Hook
log.AddHook(&MyHook{})
```

## 在 HTTP 中间件中使用

```go
func LoggingMiddleware(log *logger.Log) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // 记录请求开始
            log.WithFields(logger.Fields{
                "method":     r.Method,
                "path":       r.URL.Path,
                "remote_ip":  r.RemoteAddr,
                "user_agent": r.UserAgent(),
            }).Info("请求开始")
            
            // 执行下一个中间件
            next.ServeHTTP(w, r)
            
            // 记录请求结束
            log.WithFields(logger.Fields{
                "method":   r.Method,
                "path":     r.URL.Path,
                "duration": time.Since(start),
            }).Info("请求结束")
        })
    }
}
```

## 最佳实践

1. **结构化日志** - 使用 JSON 格式便于日志分析
2. **适当的日志级别** - 根据重要性选择合适的日志级别
3. **上下文信息** - 记录足够的上下文信息便于问题排查
4. **性能考虑** - 避免在高频路径记录过多调试日志
5. **敏感信息** - 不要记录密码、token 等敏感信息

```go
// 好的日志实践
log.WithFields(logger.Fields{
    "user_id":    userID,
    "action":     "update_profile",
    "ip_address": clientIP,
    "duration":   time.Since(start),
    "success":    true,
}).Info("用户操作完成")

// 避免记录敏感信息
log.WithFields(logger.Fields{
    "username": user.Username,
    // "password": user.Password,  // ❌ 不要记录密码
    // "token":    user.Token,     // ❌ 不要记录 token
}).Info("用户登录")
```

## 日志级别使用指南

- **Debug**: 详细的调试信息，仅在开发和调试时使用
- **Info**: 一般性信息，记录程序的正常运行状态
- **Warn**: 警告信息，程序可以继续运行但需要注意
- **Error**: 错误信息，程序遇到错误但仍可继续运行
- **Fatal**: 严重错误，程序无法继续运行，会调用 `os.Exit(1)`
- **Panic**: 恐慌错误，会触发 panic，程序崩溃