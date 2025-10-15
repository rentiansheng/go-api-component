# Decode 数据解码器

负责 HTTP 请求数据的解析和解码，支持多种数据格式。

## 支持的数据格式

- **Query String** - URL 查询参数
- **Form Data** - `application/x-www-form-urlencoded`  
- **JSON** - `application/json`
- **Multipart Form** - `multipart/form-data`

## 文件说明

- `form.go` - Form 数据解析器
- `query.go` - Query String 解析器

## 使用方式

解码器会自动根据请求的 Content-Type 选择合适的解析方式：

```go
// 自动解码（推荐）
var data RequestData
err := ctx.Decode(&data)

// 强制 JSON 解码
var data RequestData  
err := ctx.JSONDecode(&data)
```

## 解码优先级

1. Query String 参数（总是解析）
2. 根据 Content-Type 解析请求体：
   - `application/json` → JSON 解析
   - `application/x-www-form-urlencoded` → Form 解析
   - `multipart/form-data` → Multipart 解析
   - 默认 → JSON 解析

## 数据合并

当同时存在 Query String 和请求体数据时，会按以下规则合并：

1. 先解析 Query String
2. 再解析请求体数据
3. 请求体数据会覆盖同名的 Query String 参数