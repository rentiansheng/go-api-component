# Go API Component

一个功能丰富的 Go API 开发组件库，提供了构建 RESTful API 所需的常用功能和中间件。

## 特性

- 🚀 基于 `go-restful` 的路由系统
- 🛡️ 内置错误处理和错误码管理
- 📝 集成日志系统（基于 logrus）
- 🔧 灵活的配置管理（基于 viper）
- 🌐 CORS 支持
- 📦 请求/响应数据解析和验证
- 🔄 中间件支持
- 📊 完整的测试覆盖

## 依赖

- Go 1.24.5+
- github.com/emicklei/go-restful/v3
- github.com/go-playground/validator/v10
- github.com/sirupsen/logrus
- github.com/spf13/viper

## 安装

```bash
go get github.com/rentiansheng/go-api-component
```

## 快速开始

```go
package main

import (
    "github.com/rentiansheng/go-api-component/server"
    "github.com/rentiansheng/go-api-component/middleware"
    "github.com/rentiansheng/go-api-component/pkg/logger"
)

func main() {
    // 初始化日志
    log := logger.NewLogger()
    
    // 服务器配置
    serverConfig := server.Server{
        Port:         ":8080",
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
    }
    
    // 创建 HTTP 服务器
    srv := server.NewHttpServer("api-server", serverConfig, log)
    
    // 启动服务器
    if err := srv.Run(); err != nil {
        log.Fatal(err)
    }
}
```

## 项目结构

```
├── go.mod              # Go 模块文件
├── logger/             # 日志组件
├── middleware/         # 中间件组件
│   ├── context/        # 请求上下文处理
│   ├── errors/         # 错误处理
│   └── ...
├── pkg/                # 公共包
│   ├── config/         # 配置管理
│   └── logger/         # 日志工具
└── server/             # HTTP 服务器
    └── router/         # 路由管理
```

## 核心组件

### 中间件 (Middleware)

提供了完整的 HTTP 中间件支持，包括：
- 路由管理
- 请求上下文
- 错误处理
- 数据解析和验证

### 服务器 (Server)

基于 `go-restful` 的 HTTP 服务器实现，支持：
- 优雅关闭
- CORS 配置
- 超时设置
- 自定义路由

### 配置管理 (Config)

基于 `viper` 的配置管理系统，支持：
- 多种配置文件格式
- 环境变量绑定
- 默认值设置

### 日志系统 (Logger)

基于 `logrus` 的日志系统，支持：
- 结构化日志
- 日志轮转
- 多种输出格式

## 使用示例

### 创建 Web 服务

```go
// 创建 Web 实例
web := middleware.NewWeb("/api/v1")

// 定义路由
web.Get("/users").Handler(getUsersHandler)
web.Post("/users").Handler(createUserHandler)
web.Put("/users/{id}").Handler(updateUserHandler)
web.Delete("/users/{id}").Handler(deleteUserHandler)

// 获取路由服务
routes := web.Routes()
```

### 请求数据解析

```go
func createUserHandler(ctx middleware.Contexts) (interface{}, middleware.Error) {
    var user User
    if err := ctx.JSONDecode(&user); err != nil {
        return nil, err
    }
    
    // 处理业务逻辑
    result, err := createUser(user)
    if err != nil {
        return nil, middleware.NewError(500, "创建用户失败")
    }
    
    return result, nil
}
```

### 错误处理

```go
import "github.com/rentiansheng/go-api-component/middleware/errors"

// 返回业务错误
return nil, errors.NewError(400, "参数错误")

// 返回带错误码的错误
return nil, errors.NewErrorWithCode(errors.InvalidParam, "用户名不能为空")
```

## 测试

运行所有测试：

```bash
go test ./...
```

运行特定包的测试：

```bash
go test ./middleware/context
go test ./middleware/errors
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

本项目采用 MIT 许可证。

## 联系方式

- 作者: rentiansheng
- GitHub: [go-api-component](https://github.com/rentiansheng/go-api-component)