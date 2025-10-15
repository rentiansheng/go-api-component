# Middleware 中间件组件

提供 HTTP 中间件功能，包括路由管理、请求上下文处理、错误处理等核心功能。

## 组件结构

```
middleware/
├── option.go               # 中间件选项配置
├── resource_wrapper.go     # 资源包装器
├── route.go               # 路由管理
├── context/               # 请求上下文处理
├── errors/                # 错误处理
└── README.md
```

## 主要功能

### 路由管理 (Route)

提供 RESTful API 路由管理功能：

```go
// 创建 Web 实例
web := middleware.NewWeb("/api/v1")

// 定义路由
web.Get("/users").Handler(getUsersHandler)
web.Post("/users").Handler(createUserHandler)
web.Put("/users/{id}").Handler(updateUserHandler)
web.Delete("/users/{id}").Handler(deleteUserHandler)
```

### 支持的 HTTP 方法

- `GET` - 获取资源
- `POST` - 创建资源
- `PUT` - 更新资源
- `DELETE` - 删除资源
- `PATCH` - 部分更新资源
- `HEAD` - 获取资源头信息
- `OPTIONS` - 获取支持的方法

### 路由特性

- 路径参数支持：`/users/{id}`
- 查询参数支持
- 中间件链支持
- 登录验证控制

```go
// 需要登录的路由
web.Get("/profile").NeedLogin().Handler(getProfileHandler)

// 无需登录的路由
web.Get("/public").NoLogin().Handler(getPublicHandler)
```

### 内容类型支持

支持多种内容类型：

- `ContentTypeJSON` - application/json
- `ContentTypeXML` - application/xml
- `ContentTypeProtoBuf` - application/x-protobuf
- `ContentTypeMsgPack` - application/x-msgpack
- `ContentTypeYaml` - application/x-yaml
- `ContentTypeToml` - application/toml

## 使用示例

### 创建完整的 API

```go
package main

import (
    "github.com/rentiansheng/go-api-component/middleware"
)

func main() {
    // 创建 Web 服务
    web := middleware.NewWeb("/api/v1")
    
    // 用户管理路由
    web.Get("/users").Handler(listUsers)
    web.Get("/users/{id}").Handler(getUser)
    web.Post("/users").Handler(createUser)
    web.Put("/users/{id}").NeedLogin().Handler(updateUser)
    web.Delete("/users/{id}").NeedLogin().Handler(deleteUser)
    
    // 获取 restful WebService
    service := web.Routes()
    
    // 添加到容器
    container := restful.NewContainer()
    container.Add(service)
}

// 处理函数示例
func listUsers(ctx middleware.Contexts) (interface{}, middleware.Error) {
    // 获取查询参数
    page := ctx.Query("page")
    limit := ctx.Query("limit")
    
    // 业务逻辑处理
    users, err := getUserList(page, limit)
    if err != nil {
        return nil, middleware.NewError(500, "获取用户列表失败")
    }
    
    return users, nil
}

func getUser(ctx middleware.Contexts) (interface{}, middleware.Error) {
    // 获取路径参数
    userID := ctx.PathParameter("id")
    
    // 业务逻辑处理
    user, err := getUserByID(userID)
    if err != nil {
        return nil, middleware.NewError(404, "用户不存在")
    }
    
    return user, nil
}
```

## 子组件

- [Context 请求上下文](./context/README.md) - 请求数据解析、响应处理
- [Errors 错误处理](./errors/README.md) - 统一错误处理和错误码管理