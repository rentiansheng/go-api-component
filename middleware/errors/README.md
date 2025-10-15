# Errors 错误处理组件

提供统一的错误处理机制，包括错误码管理、错误消息国际化等功能。

## 组件结构

```
errors/
├── error.go              # 错误接口和基础实现
├── code/                 # 错误码定义
│   └── common.go         # 通用错误码
├── message/              # 错误消息管理
│   ├── init.og.go        # 消息初始化（注：文件名可能有误）
│   └── default/          # 默认消息
│       ├── association.go # 关联消息
│       └── lang.go       # 语言支持
└── register/             # 错误注册器
    └── register.go       # 错误注册逻辑
```

## 核心功能

### 1. 错误接口定义

```go
type Error interface {
    Error() string
    Code() int
    Message() string
    Details() interface{}
}
```

### 2. 创建错误

#### 基础错误创建
```go
import "github.com/rentiansheng/go-api-component/middleware/errors"

// 创建带状态码和消息的错误
err := errors.NewError(400, "参数错误")

// 创建带详细信息的错误
err := errors.NewErrorWithDetails(400, "验证失败", map[string]string{
    "name":  "name is required",
    "email": "email format invalid",
})
```

#### 使用错误码创建
```go
import (
    "github.com/rentiansheng/go-api-component/middleware/errors"
    "github.com/rentiansheng/go-api-component/middleware/errors/code"
)

// 使用预定义错误码
err := errors.NewErrorWithCode(code.InvalidParam, "用户名格式错误")
err := errors.NewErrorWithCode(code.NotFound, "用户不存在")
err := errors.NewErrorWithCode(code.InternalError, "服务器内部错误")
```

### 3. 常用错误码

```go
// 通用错误码（在 code/common.go 中定义）
const (
    Success      = 0     // 成功
    InvalidParam = 400   // 参数错误
    Unauthorized = 401   // 未授权
    Forbidden    = 403   // 禁止访问
    NotFound     = 404   // 资源不存在
    MethodNotAllowed = 405 // 方法不允许
    InternalError = 500  // 服务器内部错误
    BadGateway   = 502   // 网关错误
    ServiceUnavailable = 503 // 服务不可用
)
```

## 使用示例

### 在 Handler 中使用

```go
func getUserByID(ctx middleware.Contexts) (interface{}, middleware.Error) {
    // 获取路径参数
    userID := ctx.PathParameter("id")
    if userID == "" {
        return nil, errors.NewError(400, "用户ID不能为空")
    }
    
    // 验证 ID 格式
    if _, err := strconv.Atoi(userID); err != nil {
        return nil, errors.NewErrorWithCode(code.InvalidParam, "用户ID格式错误")
    }
    
    // 查询用户
    user, err := findUserByID(userID)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.NewErrorWithCode(code.NotFound, "用户不存在")
        }
        return nil, errors.NewErrorWithCode(code.InternalError, "查询用户失败")
    }
    
    return user, nil
}
```

### 数据验证错误

```go
func createUser(ctx middleware.Contexts) (interface{}, middleware.Error) {
    var req CreateUserRequest
    if err := ctx.JSONDecode(&req); err != nil {
        // JSONDecode 会自动返回验证错误
        return nil, err
    }
    
    // 业务逻辑验证
    if existsUser(req.Email) {
        return nil, errors.NewError(409, "邮箱已存在")
    }
    
    user, err := createUserInDB(req)
    if err != nil {
        return nil, errors.NewErrorWithCode(code.InternalError, "创建用户失败")
    }
    
    return user, nil
}
```

### 批量验证错误

```go
func validateUser(user User) middleware.Error {
    var details []ValidationError
    
    if user.Name == "" {
        details = append(details, ValidationError{
            Field:   "name",
            Message: "姓名不能为空",
        })
    }
    
    if !isValidEmail(user.Email) {
        details = append(details, ValidationError{
            Field:   "email", 
            Message: "邮箱格式错误",
        })
    }
    
    if len(details) > 0 {
        return errors.NewErrorWithDetails(400, "参数验证失败", details)
    }
    
    return nil
}
```

## 错误响应格式

标准错误响应格式：

```json
{
    "code": 400,
    "message": "参数验证失败",
    "data": null,
    "details": [
        {
            "field": "name",
            "message": "姓名不能为空"
        },
        {
            "field": "email",
            "message": "邮箱格式错误"
        }
    ]
}
```

## 错误码分类

### HTTP 状态码对应

- **2xx** - 成功状态
  - `200` - 操作成功
  
- **4xx** - 客户端错误
  - `400` - 参数错误
  - `401` - 未授权
  - `403` - 禁止访问
  - `404` - 资源不存在
  - `409` - 资源冲突
  
- **5xx** - 服务器错误
  - `500` - 内部服务器错误
  - `502` - 网关错误
  - `503` - 服务不可用

### 业务错误码

建议使用 5 位数字业务错误码：

```go
const (
    // 用户相关错误 10xxx
    UserNotFound     = 10001
    UserExists       = 10002
    UserDisabled     = 10003
    
    // 权限相关错误 20xxx
    PermissionDenied = 20001
    TokenExpired     = 20002
    TokenInvalid     = 20003
    
    // 业务逻辑错误 30xxx
    OrderNotFound    = 30001
    OrderCancelled   = 30002
    InsufficientBalance = 30003
)
```

## 国际化支持

错误消息支持多语言：

```go
// 注册错误消息
errors.RegisterMessage("zh-CN", code.UserNotFound, "用户不存在")
errors.RegisterMessage("en-US", code.UserNotFound, "User not found")

// 使用时会根据请求头的 Accept-Language 自动选择语言
err := errors.NewErrorWithCode(code.UserNotFound)
```

## 最佳实践

1. **统一错误码** - 为每种错误类型定义唯一的错误码
2. **分层错误处理** - 在不同层级适当转换错误类型
3. **详细错误信息** - 为开发和调试提供足够的错误详情
4. **用户友好** - 对外暴露的错误信息要用户友好
5. **日志记录** - 重要错误要记录到日志系统

```go
// 好的实践
func (s *UserService) GetUser(id string) (*User, error) {
    user, err := s.repo.GetUser(id)
    if err != nil {
        if err == sql.ErrNoRows {
            // 转换为业务错误
            return nil, errors.NewErrorWithCode(code.UserNotFound, "用户不存在")
        }
        // 记录详细错误日志
        s.logger.Error("获取用户失败", "userID", id, "error", err)
        return nil, errors.NewErrorWithCode(code.InternalError, "服务暂时不可用")
    }
    return user, nil
}
```