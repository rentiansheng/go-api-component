# Context 请求上下文处理

提供 HTTP 请求上下文处理功能，包括请求数据解析、参数获取、响应处理等。

## 核心接口

### Contexts 接口

主要提供以下功能：

```go
type Contexts interface {
    // 响应处理
    Response() http.ResponseWriter
    
    // 数据解析
    JSONDecode(target interface{}) errors.Error
    Decode(target interface{}) errors.Error
    
    // 文件处理
    FromFile(name string) (multipart.File, *multipart.FileHeader, error)
    
    // 参数获取
    Query(name string) []string
    PathParameter(name string) string
    PathParameters() map[string]string
    
    // 响应设置
    SetExtraResponse(key string, val interface{})
    SetPageResponse(val interface{})
    SetResponseFile(fileName string, content *bytes.Buffer)
    SetRawResponse(typ string, body []byte)
}
```

## 主要功能

### 1. 数据解析

#### JSON 解析
```go
func createUser(ctx middleware.Contexts) (interface{}, middleware.Error) {
    var user User
    if err := ctx.JSONDecode(&user); err != nil {
        return nil, err
    }
    
    // 处理用户创建逻辑
    return user, nil
}
```

#### 自动解析（支持 JSON、Form、Query）
```go
func updateUser(ctx middleware.Contexts) (interface{}, middleware.Error) {
    var req UpdateUserRequest
    if err := ctx.Decode(&req); err != nil {
        return nil, err
    }
    
    // 处理用户更新逻辑
    return req, nil
}
```

### 2. 参数获取

#### 路径参数
```go
// 路由定义：GET /users/{id}
userID := ctx.PathParameter("id")

// 获取所有路径参数
params := ctx.PathParameters()
```

#### 查询参数
```go
// GET /users?page=1&limit=10&tags=golang&tags=web
page := ctx.Query("page")      // ["1"]
limit := ctx.Query("limit")    // ["10"]
tags := ctx.Query("tags")      // ["golang", "web"]
```

### 3. 文件上传处理

```go
func uploadFile(ctx middleware.Contexts) (interface{}, middleware.Error) {
    file, header, err := ctx.FromFile("upload")
    if err != nil {
        return nil, middleware.NewError(400, "文件上传失败")
    }
    defer file.Close()
    
    // 处理文件
    content, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, middleware.NewError(500, "读取文件失败")
    }
    
    return map[string]interface{}{
        "filename": header.Filename,
        "size":     len(content),
    }, nil
}
```

### 4. 响应处理

#### 设置额外响应数据
```go
func getUserList(ctx middleware.Contexts) (interface{}, middleware.Error) {
    users, total, err := getUsersFromDB()
    if err != nil {
        return nil, middleware.NewError(500, "查询失败")
    }
    
    // 设置分页信息
    ctx.SetPageResponse(map[string]interface{}{
        "total": total,
        "page":  1,
        "limit": 10,
    })
    
    // 设置额外信息
    ctx.SetExtraResponse("timestamp", time.Now().Unix())
    
    return users, nil
}
```

#### 文件下载响应
```go
func downloadFile(ctx middleware.Contexts) (interface{}, middleware.Error) {
    content := generateFileContent()
    buffer := bytes.NewBuffer(content)
    
    ctx.SetResponseFile("report.pdf", buffer)
    return nil, nil
}
```

#### 原始响应
```go
func customResponse(ctx middleware.Contexts) (interface{}, middleware.Error) {
    xmlData := `<?xml version="1.0"?><root><message>Hello</message></root>`
    ctx.SetRawResponse("application/xml", []byte(xmlData))
    return nil, nil
}
```

## 数据验证

支持使用 `validator` 标签进行数据验证：

```go
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"min=0,max=150"`
    Phone    string `json:"phone" validate:"omitempty,len=11"`
}

func createUser(ctx middleware.Contexts) (interface{}, middleware.Error) {
    var req CreateUserRequest
    if err := ctx.JSONDecode(&req); err != nil {
        // 自动返回验证错误信息
        return nil, err
    }
    
    // 数据已通过验证，可以安全使用
    return createUserInDB(req), nil
}
```

## 验证规则

支持的验证规则包括：
- `required` - 必填
- `min`, `max` - 最小/最大值（数字）或长度（字符串）
- `len` - 固定长度
- `email` - 邮箱格式
- `url` - URL 格式
- `alpha` - 只允许字母
- `alphanum` - 只允许字母和数字
- `numeric` - 只允许数字
- `oneof` - 枚举值验证

更多验证规则请参考：[validator 文档](https://github.com/go-playground/validator#baked-in-validations)

## 解码支持

### 支持的数据格式

1. **Query String** - URL 查询参数
2. **Form Data** - `application/x-www-form-urlencoded`
3. **JSON** - `application/json`
4. **Multipart Form** - `multipart/form-data`

### 自动解码流程

1. 首先解析 Query String 参数
2. 根据 Content-Type 解析请求体：
   - `application/json` → JSON 解析
   - `application/x-www-form-urlencoded` → Form 解析
   - `multipart/form-data` → Multipart 解析
3. 合并所有解析的数据
4. 执行数据验证

## 错误处理

当解析或验证失败时，会返回详细的错误信息：

```go
// 验证失败示例响应
{
    "code": 400,
    "message": "参数验证失败",
    "details": [
        {
            "field": "name",
            "message": "name is required"
        },
        {
            "field": "email", 
            "message": "email must be a valid email address"
        }
    ]
}
```

## 测试

运行上下文相关的测试：

```bash
go test ./middleware/context
```