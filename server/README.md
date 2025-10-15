# Server HTTP 服务器

基于 `go-restful` 的 HTTP 服务器实现，提供完整的 Web 服务功能。

## 组件结构

```
server/
├── server.go           # HTTP 服务器核心实现
├── router/             # 路由管理
│   ├── router.go       # 路由器实现
│   └── README.md
└── README.md
```

## 核心功能

- 🚀 基于 `go-restful` 的 RESTful API 服务器
- 🛡️ 优雅关闭支持
- 🌐 CORS 跨域支持
- ⏱️ 可配置的超时设置
- 📝 集成日志系统
- 🔧 灵活的服务器配置

## 服务器配置

```go
type Server struct {
    Port            string        `mapstructure:"port"`             // 服务端口
    ReadTimeout     time.Duration `mapstructure:"read_timeout"`     // 读取超时
    WriteTimeout    time.Duration `mapstructure:"write_timeout"`    // 写入超时
    ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"` // 关闭超时
    Cors            Cors          `mapstructure:"cors"`             // CORS 配置
}

type Cors struct {
    AllowedHeaders []string `mapstructure:"allowed_headers"` // 允许的请求头
    AllowedMethods []string `mapstructure:"allowed_methods"` // 允许的 HTTP 方法
    AllowedDomains []string `mapstructure:"allowed_domains"` // 允许的域名
    CookiesAllowed bool     `mapstructure:"cookies_allowed"` // 是否允许 cookies
    EnableCORS     bool     `mapstructure:"enable_cors"`     // 是否启用 CORS
}
```

## 基础使用

```go
package main

import (
    "time"
    
    "github.com/rentiansheng/go-api-component/server"
    "github.com/rentiansheng/go-api-component/pkg/logger"
)

func main() {
    // 初始化日志
    log := logger.NewLogger()
    
    // 服务器配置
    serverConfig := server.Server{
        Port:            ":8080",
        ReadTimeout:     30 * time.Second,
        WriteTimeout:    30 * time.Second,
        ShutdownTimeout: 10 * time.Second,
        Cors: server.Cors{
            EnableCORS:     true,
            AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
            AllowedHeaders: []string{"Content-Type", "Authorization"},
            AllowedDomains: []string{"*"},
            CookiesAllowed: false,
        },
    }
    
    // 创建 HTTP 服务器
    srv := server.NewHttpServer("api-server", serverConfig, log)
    
    // 启动服务器
    if err := srv.Run(); err != nil {
        log.Fatal("服务器启动失败:", err)
    }
}
```

## 完整应用示例

```go
package main

import (
    "time"
    
    "github.com/emicklei/go-restful/v3"
    "github.com/rentiansheng/go-api-component/server"
    "github.com/rentiansheng/go-api-component/server/router"
    "github.com/rentiansheng/go-api-component/middleware"
    "github.com/rentiansheng/go-api-component/pkg/logger"
)

type UserHandler struct {
    log *logger.Log
}

func (h *UserHandler) listUsers(ctx middleware.Contexts) (interface{}, middleware.Error) {
    // 模拟用户列表
    users := []map[string]interface{}{
        {"id": 1, "name": "John Doe", "email": "john@example.com"},
        {"id": 2, "name": "Jane Smith", "email": "jane@example.com"},
    }
    
    h.log.Info("获取用户列表")
    return users, nil
}

func (h *UserHandler) getUser(ctx middleware.Contexts) (interface{}, middleware.Error) {
    userID := ctx.PathParameter("id")
    
    // 模拟用户查询
    user := map[string]interface{}{
        "id":    userID,
        "name":  "John Doe",
        "email": "john@example.com",
    }
    
    h.log.WithField("userID", userID).Info("获取用户详情")
    return user, nil
}

func main() {
    // 初始化日志
    log := logger.NewLogger()
    log.SetLevel(logger.InfoLevel)
    
    // 创建处理器
    userHandler := &UserHandler{log: log}
    
    // 创建路由
    r := router.NewRouter()
    
    // 创建 Web 服务
    web := middleware.NewWeb("/api/v1")
    
    // 定义用户相关路由
    web.Get("/users").Handler(userHandler.listUsers)
    web.Get("/users/{id}").Handler(userHandler.getUser)
    
    // 注册路由到路由器
    r.RegisterRoutes(web.Routes())
    
    // 服务器配置
    serverConfig := server.Server{
        Port:            ":8080",
        ReadTimeout:     30 * time.Second,
        WriteTimeout:    30 * time.Second,
        ShutdownTimeout: 10 * time.Second,
        Cors: server.Cors{
            EnableCORS:     true,
            AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
            AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
            AllowedDomains: []string{"http://localhost:3000", "https://example.com"},
            CookiesAllowed: false,
        },
    }
    
    // 创建并启动服务器
    srv := server.NewHttpServer("user-api", serverConfig, log)
    
    log.Info("启动服务器 http://localhost:8080")
    if err := srv.Run(); err != nil {
        log.Fatal("服务器启动失败:", err)
    }
}
```

## 配置文件示例

### YAML 配置

```yaml
server:
  port: ":8080"
  read_timeout: "30s"
  write_timeout: "30s"  
  shutdown_timeout: "10s"
  cors:
    enable_cors: true
    allowed_methods:
      - "GET"
      - "POST"
      - "PUT"
      - "DELETE"
      - "OPTIONS"
    allowed_headers:
      - "Content-Type"
      - "Authorization"
      - "X-Requested-With"
    allowed_domains:
      - "http://localhost:3000"
      - "https://example.com"
    cookies_allowed: false
```

### 从配置文件加载

```go
import (
    "github.com/rentiansheng/go-api-component/pkg/config"
    "github.com/rentiansheng/go-api-component/server"
)

func main() {
    // 加载配置
    cfg := config.NewConfigHandlerWithDefaults("app")
    if err := cfg.Load(); err != nil {
        log.Fatal("加载配置失败:", err)
    }
    
    // 解析服务器配置
    var serverConfig server.Server
    if err := cfg.UnmarshalKey("server", &serverConfig); err != nil {
        log.Fatal("解析服务器配置失败:", err)
    }
    
    // 创建服务器
    srv := server.NewHttpServer("api-server", serverConfig, log)
    srv.Run()
}
```

## CORS 配置详解

### 基础 CORS 配置

```go
cors := server.Cors{
    EnableCORS:     true,                    // 启用 CORS
    AllowedOrigins: []string{"*"},           // 允许所有域名
    AllowedMethods: []string{                // 允许的 HTTP 方法
        "GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH",
    },
    AllowedHeaders: []string{                // 允许的请求头
        "Content-Type", 
        "Authorization", 
        "X-Requested-With",
        "X-API-Key",
    },
    CookiesAllowed: false,                   // 不允许携带 cookies
}
```

### 生产环境 CORS 配置

```go
cors := server.Cors{
    EnableCORS: true,
    AllowedDomains: []string{                // 明确指定允许的域名
        "https://app.example.com",
        "https://admin.example.com", 
        "https://mobile.example.com",
    },
    AllowedMethods: []string{                // 限制允许的方法
        "GET", "POST", "PUT", "DELETE",
    },
    AllowedHeaders: []string{                // 限制允许的请求头
        "Content-Type",
        "Authorization",
    },
    CookiesAllowed: true,                    // 允许携带 cookies（如需要）
}
```

## 优雅关闭

服务器支持优雅关闭，会等待现有请求处理完成：

```go
// 监听系统信号
go func() {
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    <-sigCh
    
    log.Info("收到关闭信号，开始优雅关闭服务器")
    
    // 创建关闭上下文
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // 关闭服务器
    if err := srv.Shutdown(ctx); err != nil {
        log.Error("服务器关闭失败:", err)
    } else {
        log.Info("服务器已优雅关闭")
    }
}()
```

## 中间件集成

```go
import (
    "github.com/emicklei/go-restful/v3"
)

// 日志中间件
func loggingFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
    start := time.Now()
    
    log.WithFields(logger.Fields{
        "method": req.Request.Method,
        "path":   req.Request.URL.Path,
        "ip":     req.Request.RemoteAddr,
    }).Info("请求开始")
    
    chain.ProcessFilter(req, resp)
    
    log.WithFields(logger.Fields{
        "method":   req.Request.Method,
        "path":     req.Request.URL.Path,
        "status":   resp.StatusCode(),
        "duration": time.Since(start),
    }).Info("请求完成")
}

// 恢复中间件
func recoveryFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
    defer func() {
        if err := recover(); err != nil {
            log.WithFields(logger.Fields{
                "panic": err,
                "path":  req.Request.URL.Path,
            }).Error("处理请求时发生 panic")
            
            resp.WriteErrorString(500, "Internal Server Error")
        }
    }()
    
    chain.ProcessFilter(req, resp)
}

// 添加中间件
container := restful.NewContainer()
container.Filter(loggingFilter)
container.Filter(recoveryFilter)
```

## 性能优化

### 连接池配置

```go
import "net/http"

func optimizedServer(config server.Server) *http.Server {
    return &http.Server{
        Addr:         config.Port,
        ReadTimeout:  config.ReadTimeout,
        WriteTimeout: config.WriteTimeout,
        IdleTimeout:  120 * time.Second,        // 空闲连接超时
        MaxHeaderBytes: 1 << 20,                // 1MB 请求头限制
    }
}
```

### 资源限制

```go
// 限制并发连接数
func limitConnections(max int) func(http.Handler) http.Handler {
    sem := make(chan struct{}, max)
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            sem <- struct{}{}
            defer func() { <-sem }()
            next.ServeHTTP(w, r)
        })
    }
}
```

## 监控和健康检查

```go
// 健康检查端点
func healthCheck(ctx middleware.Contexts) (interface{}, middleware.Error) {
    return map[string]interface{}{
        "status":    "healthy",
        "timestamp": time.Now().Unix(),
        "version":   "1.0.0",
    }, nil
}

// 指标端点
func metrics(ctx middleware.Contexts) (interface{}, middleware.Error) {
    return map[string]interface{}{
        "requests_total": getTotalRequests(),
        "uptime_seconds": getUptimeSeconds(),
        "memory_usage":   getMemoryUsage(),
    }, nil
}

// 注册监控端点
web.Get("/health").NoLogin().Handler(healthCheck)
web.Get("/metrics").NoLogin().Handler(metrics)
```

## 最佳实践

1. **合理的超时设置** - 根据业务需求设置合适的超时时间
2. **CORS 安全配置** - 生产环境不要使用 `*` 通配符
3. **优雅关闭** - 确保服务器能够优雅关闭，处理完现有请求
4. **监控集成** - 添加健康检查和指标端点
5. **错误处理** - 统一的错误处理和日志记录
6. **安全考虑** - 添加必要的安全中间件（如速率限制、认证等）

详细的路由管理请参考：[Router 路由管理](./router/README.md)