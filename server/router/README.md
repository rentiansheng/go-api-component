# Router 路由管理

提供 RESTful API 路由管理和注册功能。

## 功能特性

- 🎯 路由注册和管理
- 🔗 与 `go-restful` 无缝集成
- 📝 支持路由组和中间件
- 🛡️ 统一的错误处理
- 📊 路由统计和监控

## 基础使用

```go
import (
    "github.com/rentiansheng/go-api-component/server/router"
    "github.com/rentiansheng/go-api-component/middleware"
)

// 创建路由器
r := router.NewRouter()

// 创建 Web 服务
web := middleware.NewWeb("/api/v1")

// 定义路由
web.Get("/users").Handler(getUsersHandler)
web.Post("/users").Handler(createUserHandler)
web.Put("/users/{id}").Handler(updateUserHandler)
web.Delete("/users/{id}").Handler(deleteUserHandler)

// 注册路由
r.RegisterRoutes(web.Routes())
```

## 路由组织

### 按功能模块组织

```go
// 用户模块
userWeb := middleware.NewWeb("/api/v1/users")
userWeb.Get("").Handler(listUsers)
userWeb.Get("/{id}").Handler(getUser)
userWeb.Post("").Handler(createUser)
userWeb.Put("/{id}").Handler(updateUser)
userWeb.Delete("/{id}").Handler(deleteUser)

// 订单模块  
orderWeb := middleware.NewWeb("/api/v1/orders")
orderWeb.Get("").Handler(listOrders)
orderWeb.Get("/{id}").Handler(getOrder)
orderWeb.Post("").Handler(createOrder)

// 注册所有路由
r.RegisterRoutes(userWeb.Routes())
r.RegisterRoutes(orderWeb.Routes())
```

### 按版本组织

```go
// API v1
v1Web := middleware.NewWeb("/api/v1")
v1Web.Get("/users").Handler(v1GetUsers)
v1Web.Get("/orders").Handler(v1GetOrders)

// API v2
v2Web := middleware.NewWeb("/api/v2")  
v2Web.Get("/users").Handler(v2GetUsers)
v2Web.Get("/orders").Handler(v2GetOrders)

// 注册不同版本的路由
r.RegisterRoutes(v1Web.Routes())
r.RegisterRoutes(v2Web.Routes())
```

## 路由示例

### RESTful 资源路由

```go
// 用户资源的完整 RESTful 路由
func setupUserRoutes() *restful.WebService {
    web := middleware.NewWeb("/api/v1/users")
    
    // GET /api/v1/users - 获取用户列表
    web.Get("").Handler(func(ctx middleware.Contexts) (interface{}, middleware.Error) {
        // 获取查询参数
        page := ctx.Query("page")
        limit := ctx.Query("limit")
        search := ctx.Query("search")
        
        users, total, err := getUserList(page, limit, search)
        if err != nil {
            return nil, middleware.NewError(500, "获取用户列表失败")
        }
        
        // 设置分页信息
        ctx.SetPageResponse(map[string]interface{}{
            "total": total,
            "page":  page,
            "limit": limit,
        })
        
        return users, nil
    })
    
    // GET /api/v1/users/{id} - 获取单个用户
    web.Get("/{id}").Handler(func(ctx middleware.Contexts) (interface{}, middleware.Error) {
        userID := ctx.PathParameter("id")
        
        user, err := getUserByID(userID)
        if err != nil {
            return nil, middleware.NewError(404, "用户不存在")
        }
        
        return user, nil
    })
    
    // POST /api/v1/users - 创建用户
    web.Post("").Handler(func(ctx middleware.Contexts) (interface{}, middleware.Error) {
        var req CreateUserRequest
        if err := ctx.JSONDecode(&req); err != nil {
            return nil, err
        }
        
        user, err := createUser(req)
        if err != nil {
            return nil, middleware.NewError(400, "创建用户失败")
        }
        
        return user, nil
    })
    
    // PUT /api/v1/users/{id} - 更新用户
    web.Put("/{id}").NeedLogin().Handler(func(ctx middleware.Contexts) (interface{}, middleware.Error) {
        userID := ctx.PathParameter("id")
        
        var req UpdateUserRequest
        if err := ctx.JSONDecode(&req); err != nil {
            return nil, err
        }
        
        user, err := updateUser(userID, req)
        if err != nil {
            return nil, middleware.NewError(400, "更新用户失败")
        }
        
        return user, nil
    })
    
    // DELETE /api/v1/users/{id} - 删除用户
    web.Delete("/{id}").NeedLogin().Handler(func(ctx middleware.Contexts) (interface{}, middleware.Error) {
        userID := ctx.PathParameter("id")
        
        if err := deleteUser(userID); err != nil {
            return nil, middleware.NewError(400, "删除用户失败")
        }
        
        return map[string]string{"message": "删除成功"}, nil
    })
    
    return web.Routes()
}
```

### 嵌套资源路由

```go
// 用户的订单资源
func setupUserOrderRoutes() *restful.WebService {
    web := middleware.NewWeb("/api/v1/users")
    
    // GET /api/v1/users/{userId}/orders - 获取用户的订单列表
    web.Get("/{userId}/orders").Handler(func(ctx middleware.Contexts) (interface{}, middleware.Error) {
        userID := ctx.PathParameter("userId")
        
        orders, err := getUserOrders(userID)
        if err != nil {
            return nil, middleware.NewError(500, "获取用户订单失败")
        }
        
        return orders, nil
    })
    
    // GET /api/v1/users/{userId}/orders/{orderId} - 获取用户的特定订单
    web.Get("/{userId}/orders/{orderId}").Handler(func(ctx middleware.Contexts) (interface{}, middleware.Error) {
        userID := ctx.PathParameter("userId")
        orderID := ctx.PathParameter("orderId")
        
        order, err := getUserOrder(userID, orderID)
        if err != nil {
            return nil, middleware.NewError(404, "订单不存在")
        }
        
        return order, nil
    })
    
    // POST /api/v1/users/{userId}/orders - 为用户创建订单
    web.Post("/{userId}/orders").NeedLogin().Handler(func(ctx middleware.Contexts) (interface{}, middleware.Error) {
        userID := ctx.PathParameter("userId")
        
        var req CreateOrderRequest
        if err := ctx.JSONDecode(&req); err != nil {
            return nil, err
        }
        
        order, err := createUserOrder(userID, req)
        if err != nil {
            return nil, middleware.NewError(400, "创建订单失败")
        }
        
        return order, nil
    })
    
    return web.Routes()
}
```

## 路由中间件

### 认证中间件

```go
// 认证检查
func requireAuth(ctx middleware.Contexts) middleware.Error {
    token := ctx.Request().Header.Get("Authorization")
    if token == "" {
        return middleware.NewError(401, "缺少认证令牌")
    }
    
    user, err := validateToken(token)
    if err != nil {
        return middleware.NewError(401, "无效的认证令牌")
    }
    
    // 将用户信息存储到上下文
    ctx.SetValue("user", user)
    return nil
}

// 需要认证的路由
web.Get("/profile").NeedLogin().Handler(getProfileHandler)
```

### 权限中间件

```go
func requirePermission(permission string) middleware.RouteOption {
    return func(route middleware.Route) middleware.Route {
        return route.Middleware(func(ctx middleware.Contexts) middleware.Error {
            user := ctx.Value("user").(*User)
            
            if !user.HasPermission(permission) {
                return middleware.NewError(403, "权限不足")
            }
            
            return nil
        })
    }
}

// 需要特定权限的路由
web.Delete("/users/{id}").
    NeedLogin().
    Middleware(requirePermission("user:delete")).
    Handler(deleteUserHandler)
```

## 路由文档生成

```go
import "github.com/emicklei/go-restful/v3"

// 添加 API 文档信息
web.Doc("用户管理 API").
    Consumes(restful.MIME_JSON).
    Produces(restful.MIME_JSON)

// 为路由添加文档
web.Get("/users").
    Doc("获取用户列表").
    Param(web.QueryParameter("page", "页码").DataType("integer")).
    Param(web.QueryParameter("limit", "每页数量").DataType("integer")).
    Param(web.QueryParameter("search", "搜索关键词").DataType("string")).
    Returns(200, "成功", []User{}).
    Returns(500, "服务器错误", nil).
    Handler(getUsersHandler)

web.Get("/users/{id}").
    Doc("获取用户详情").
    Param(web.PathParameter("id", "用户ID").DataType("string")).
    Returns(200, "成功", User{}).
    Returns(404, "用户不存在", nil).
    Handler(getUserHandler)
```

## 路由测试

```go
func TestUserRoutes(t *testing.T) {
    // 创建测试路由器
    r := router.NewRouter()
    
    // 注册路由
    r.RegisterRoutes(setupUserRoutes())
    
    // 创建测试服务器
    server := httptest.NewServer(r.Handler())
    defer server.Close()
    
    // 测试获取用户列表
    resp, err := http.Get(server.URL + "/api/v1/users")
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    
    // 测试创建用户
    user := CreateUserRequest{
        Name:  "Test User",
        Email: "test@example.com",
    }
    
    body, _ := json.Marshal(user)
    resp, err = http.Post(
        server.URL+"/api/v1/users",
        "application/json",
        bytes.NewBuffer(body),
    )
    assert.NoError(t, err)
    assert.Equal(t, 201, resp.StatusCode)
}
```

## 性能优化

### 路由缓存

```go
type RouteCache struct {
    cache map[string]*restful.Route
    mutex sync.RWMutex
}

func (rc *RouteCache) FindRoute(method, path string) *restful.Route {
    key := method + ":" + path
    
    rc.mutex.RLock()
    route, exists := rc.cache[key]
    rc.mutex.RUnlock()
    
    if exists {
        return route
    }
    
    // 查找路由并缓存
    route = rc.findRouteFromContainer(method, path)
    if route != nil {
        rc.mutex.Lock()
        rc.cache[key] = route
        rc.mutex.Unlock()
    }
    
    return route
}
```

### 路由压缩

```go
import "github.com/klauspost/compress/gzip"

// 启用 gzip 压缩
func enableGzipCompression(container *restful.Container) {
    container.EnableContentEncoding(true)
}
```

## 最佳实践

1. **RESTful 设计** - 遵循 RESTful API 设计原则
2. **路由组织** - 按功能模块或版本合理组织路由
3. **参数验证** - 在路由层面进行基础参数验证
4. **错误处理** - 统一的错误响应格式
5. **文档化** - 为每个路由添加清晰的文档
6. **测试覆盖** - 确保所有路由都有对应的测试
7. **性能考虑** - 合理使用缓存和中间件

```go
// 路由组织示例
func SetupAPIRoutes() []*restful.WebService {
    var services []*restful.WebService
    
    // 用户相关路由
    services = append(services, setupUserRoutes())
    
    // 订单相关路由  
    services = append(services, setupOrderRoutes())
    
    // 产品相关路由
    services = append(services, setupProductRoutes())
    
    return services
}

func main() {
    r := router.NewRouter()
    
    // 注册所有路由
    for _, service := range SetupAPIRoutes() {
        r.RegisterRoutes(service)
    }
}
```