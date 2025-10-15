# Config 配置管理

基于 `viper` 的配置管理系统，提供灵活的配置加载和管理功能。

## 特性

- 🔧 支持多种配置文件格式（JSON、YAML、TOML、HCL 等）
- 🌍 环境变量绑定支持
- 🎯 默认值设置
- 📁 多目录配置文件查找
- 🔄 配置热重载
- 🎨 链式配置构建

## 使用方法

### 基础使用

```go
import "github.com/rentiansheng/go-api-component/pkg/config"

// 创建配置处理器（使用默认设置）
cfg := config.NewConfigHandlerWithDefaults("app")

// 设置配置文件查找目录
cfg.SetDirs([]string{".", "config", "/etc/myapp"})

// 加载配置
if err := cfg.Load(); err != nil {
    log.Fatal("加载配置失败:", err)
}

// 读取配置
port := cfg.GetString("server.port")
timeout := cfg.GetDuration("server.timeout")
```

### 高级配置

```go
// 自定义配置处理器
cfg := config.NewConfigHandler(
    "myapp",                    // 配置文件名
    []string{".", "config"},    // 查找目录
    map[string][]string{        // 环境变量绑定
        "server.port": {"PORT", "SERVER_PORT"},
        "database.host": {"DB_HOST"},
    },
    map[string]interface{}{     // 默认值
        "server.port": "8080",
        "server.timeout": "30s",
        "database.host": "localhost",
    },
)
```

## 配置文件示例

### YAML 格式 (app.yaml)

```yaml
server:
  port: 8080
  host: 0.0.0.0
  timeout: 30s
  cors:
    enabled: true
    origins:
      - "http://localhost:3000"
      - "https://example.com"

database:
  host: localhost
  port: 5432
  name: myapp
  username: user
  password: password

logger:
  level: info
  format: json
  output: stdout
```

### JSON 格式 (app.json)

```json
{
  "server": {
    "port": 8080,
    "host": "0.0.0.0",
    "timeout": "30s"
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "name": "myapp"
  },
  "logger": {
    "level": "info",
    "format": "json"
  }
}
```

## 环境变量绑定

```go
// 绑定环境变量
cfg.SetBindEnv(map[string][]string{
    // 配置键: [环境变量名列表]
    "server.port":    {"PORT", "SERVER_PORT"},
    "database.host":  {"DB_HOST", "DATABASE_HOST"},
    "database.port":  {"DB_PORT"},
    "logger.level":   {"LOG_LEVEL"},
})

// 环境变量优先级高于配置文件
// 例如：export PORT=9090 会覆盖配置文件中的 server.port
```

## 默认值设置

```go
// 设置默认值
cfg.SetDefault(map[string]interface{}{
    "server.port":           "8080",
    "server.host":           "0.0.0.0", 
    "server.timeout":        "30s",
    "server.read_timeout":   "15s",
    "server.write_timeout":  "15s",
    "logger.level":          "info",
    "logger.format":         "json",
    "database.max_conns":    10,
    "database.max_idle":     5,
})
```

## 配置读取方法

```go
// 字符串
port := cfg.GetString("server.port")
host := cfg.GetString("server.host")

// 数字
maxConns := cfg.GetInt("database.max_conns")
timeout := cfg.GetDuration("server.timeout")

// 布尔值
corsEnabled := cfg.GetBool("server.cors.enabled")

// 数组
origins := cfg.GetStringSlice("server.cors.origins")

// 映射
dbConfig := cfg.GetStringMapString("database")

// 复杂结构体
var serverConfig ServerConfig
if err := cfg.UnmarshalKey("server", &serverConfig); err != nil {
    log.Fatal("解析服务器配置失败:", err)
}
```

## 完整使用示例

```go
package main

import (
    "log"
    "time"
    
    "github.com/rentiansheng/go-api-component/pkg/config"
)

type AppConfig struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Logger   LoggerConfig   `mapstructure:"logger"`
}

type ServerConfig struct {
    Port         string        `mapstructure:"port"`
    Host         string        `mapstructure:"host"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Name     string `mapstructure:"name"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
}

type LoggerConfig struct {
    Level  string `mapstructure:"level"`
    Format string `mapstructure:"format"`
}

func main() {
    // 创建配置处理器
    cfg := config.NewConfigHandler(
        "app",
        []string{".", "config", "/etc/myapp"},
        map[string][]string{
            "server.port":    {"PORT"},
            "database.host":  {"DB_HOST"},
            "database.name":  {"DB_NAME"},
            "logger.level":   {"LOG_LEVEL"},
        },
        map[string]interface{}{
            "server.port":           "8080",
            "server.host":           "0.0.0.0",
            "server.read_timeout":   "30s",
            "server.write_timeout":  "30s",
            "database.host":         "localhost",
            "database.port":         5432,
            "logger.level":          "info",
            "logger.format":         "json",
        },
    )
    
    // 加载配置
    if err := cfg.Load(); err != nil {
        log.Fatal("加载配置失败:", err)
    }
    
    // 解析到结构体
    var appConfig AppConfig
    if err := cfg.Unmarshal(&appConfig); err != nil {
        log.Fatal("解析配置失败:", err)
    }
    
    // 使用配置
    log.Printf("服务器将在 %s:%s 启动", appConfig.Server.Host, appConfig.Server.Port)
    log.Printf("数据库连接: %s:%d/%s", appConfig.Database.Host, appConfig.Database.Port, appConfig.Database.Name)
}
```

## 配置优先级

配置值的优先级从高到低：

1. **显式设置** - 通过 `cfg.Set()` 显式设置
2. **命令行参数** - 通过 pflag 绑定的命令行参数
3. **环境变量** - 通过 `BindEnv` 绑定的环境变量
4. **配置文件** - 从配置文件读取的值
5. **默认值** - 通过 `SetDefault` 设置的默认值

## 最佳实践

1. **统一配置结构** - 使用结构体定义配置模式
2. **环境分离** - 不同环境使用不同的配置文件
3. **敏感信息** - 密码等敏感信息通过环境变量传入
4. **配置验证** - 加载配置后验证必要字段
5. **文档化** - 为每个配置项提供清晰的文档

```go
// 配置验证示例
func validateConfig(cfg AppConfig) error {
    if cfg.Server.Port == "" {
        return errors.New("server.port is required")
    }
    if cfg.Database.Host == "" {
        return errors.New("database.host is required")
    }
    return nil
}
```