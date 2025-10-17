package context

import (
	"bytes"
	osCtx "context"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rentiansheng/go-api-component/middleware/errors"
	_ "github.com/rentiansheng/go-api-component/middleware/errors/message"
)

const (
	spanIDKey = "span-id"
	// CtxLogIDKey context log id key
	CtxLogIDKey = "trace-id"
)

type Contexts interface {
	Response() http.ResponseWriter

	// JSONDecode 使用json格式解析request body
	// 详细的校验规则，https://github.com/go-playground/validator#baked-in-validations
	JSONDecode(target interface{}) errors.Error
	// Decode 根据http中的信息适配需要解析方式。已经支持自动解析query string, post(application/x-www-form-urlencoded,json)
	// 默认是先解析 query string ,然后根据 http content type 解析 合并数据。默认是json 方式
	// 详细的校验规则，https://github.com/go-playground/validator#baked-in-validations
	Decode(target interface{}) errors.Error
	FromFile(name string) (multipart.File, *multipart.FileHeader, error)
	Query(name string) []string
	PathParameter(name string) string
	PathParameters() map[string]string
	SetExtraResponse(key string, val interface{})
	// SetPageResponse  这里的数据回与data 同级返回给用户page字段信息，只有用contexts 返回有效
	SetPageResponse(val interface{})
	SetResponseFile(fileName string, content *bytes.Buffer)
	GetResponseFile() (fileName string, content *bytes.Buffer, exists bool)
	GetExtraResponse() map[string]interface{}

	SetRawResponse(typ string, body []byte)
	GetRawResponse() (typ string, body []byte, exists bool)

	SelectedRoutePath() string

	SetData(data interface{})
	GetData() interface{}

	HTTPBody() ([]byte, error)

	Context
}

type Context interface {
	osCtx.Context
	Cookie() []*http.Cookie
	Header() http.Header
	Request() *http.Request

	GetRequestID() string

	Log() Log
	Error() errImpl

	// IsDone 判断context 是否已经结束
	IsDone() bool
	// WithTimeout 设置超时间时间
	WithTimeout(timeout time.Duration)

	WithValue(key, value interface{})
	SubContext(suffix string) Context
	WithTimeoutCtx(timeout time.Duration) Context

	Mapper(action string, src, dst interface{}) errors.Error
	AllMapper(action string, src, dst interface{}) errors.Error

	Cancel() osCtx.CancelFunc

	WithSpan() Context
	WithSpanPrefix(prefix string) Context
	WithSpanID(id string) Context
}

type rawResponse struct {
	body []byte
	typ  string
}

func NewSysContext(ctx osCtx.Context) Context {
	if rctx, ok := ctx.(*ginContext); ok {
		newCtx := rctx.clone()
		return newCtx
	} else {
		ctx = AdjustCtxLogID(ctx)
	}
	me := &ginContext{
		requestID:     GetLogID(ctx),
		ctx:           ctx,
		meErr:         &err{},
		extraResponse: make(map[string]interface{}),
		c:             &gin.Context{},
	}
	return me
}

// NewContext 不能向 Contexts 对象转换，回panic
func NewContext(g *gin.Context) Contexts {
	ctx := g.Request.Context()
	if rctx, ok := ctx.(*ginContext); ok {
		newCtx := rctx.clone()
		return newCtx
	} else {
		ctx = AdjustCtxLogID(ctx)
	}

	me := &ginContext{
		requestID:     GetLogID(ctx),
		ctx:           ctx,
		meErr:         &err{},
		c:             g,
		extraResponse: make(map[string]interface{}),
	}
	return me
}

func TODO() Context {

	ctx := osCtx.Background()
	ctx = AdjustCtxLogID(ctx)

	me := &ginContext{
		requestID:     GetLogID(ctx),
		ctx:           ctx,
		meErr:         &err{},
		c:             &gin.Context{},
		extraResponse: make(map[string]interface{}),
	}

	return me
}
