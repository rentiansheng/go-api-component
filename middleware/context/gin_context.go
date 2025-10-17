package context

import (
	"bytes"
	osCtx "context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rentiansheng/go-api-component/middleware/errors"
	"github.com/rentiansheng/go-api-component/middleware/errors/code"
	"github.com/rentiansheng/mapper"
)

type ginContext struct {
	c             *gin.Context
	extraResponse map[string]interface{}
	pageResponse  interface{}
	data          interface{}
	requestID     string
	ctx           osCtx.Context
	responseFile  struct {
		fileName string
		content  *bytes.Buffer
		exists   bool
	}
	rawResponse struct {
		typ    string
		body   []byte
		exists bool
	}
	cancelFn func()

	meErr errImpl
}

// Cancel implements Contexts.
func (g *ginContext) Cancel() osCtx.CancelFunc {
	return g.cancelFn
}

// Cookie implements Contexts.
func (g *ginContext) Cookie() []*http.Cookie {
	return g.c.Request.Cookies()
}

// Deadline implements Contexts.
func (g *ginContext) Deadline() (deadline time.Time, ok bool) {
	return g.ctx.Deadline()
}

// Done implements Contexts.
func (g *ginContext) Done() <-chan struct{} {
	return g.ctx.Done()
}

// Err implements Contexts.
func (g *ginContext) Err() error {
	return g.ctx.Err()
}

// Error implements Contexts.
func (g *ginContext) Error() errImpl {
	return g.meErr
}

// GetData implements Contexts.
func (g *ginContext) GetData() interface{} {
	return g.data
}

// GetRequestID implements Contexts.
func (g *ginContext) GetRequestID() string {
	return g.requestID
}

// HTTPBody implements Contexts.
func (g *ginContext) HTTPBody() ([]byte, error) {
	if g.c == nil || g.c.Request == nil || g.c.Request.Body == nil {
		return nil, nil
	}

	if g.c.Request.Body == nil {
		return nil, nil
	}

	bodyBytes, err := ioutil.ReadAll(g.c.Request.Body)
	if err != nil {
		g.Log().Errorf("io read request.Body fail, %+v", err)
		return nil, g.Error().LegacyWrapCode(code.JSONDecodeErrCode, err)
	}
	g.c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	return bodyBytes, nil
}

// Header implements Contexts.
func (g *ginContext) Header() http.Header {
	if g.c == nil || g.c.Request == nil {
		return http.Header{}
	}
	return g.c.Request.Header
}

// IsDone implements Contexts.
func (g *ginContext) IsDone() bool {
	select {
	case <-g.ctx.Done():
		return true
	default:
		return false
	}
}

// Mapper implements Contexts.
func (g *ginContext) Mapper(action string, src interface{}, dst interface{}) errors.Error {
	if src == nil {
		return nil
	}
	if err := mapper.Mapper(g.c, src, dst); err != nil {
		g.Log().ErrorJSON("action: %s, src: %s, err: %s", action, src, err.Error())
		return g.Error().Errorf(code.MapperActionErrCode, action, err.Error())
	}
	return nil
}

// Request implements Contexts.
func (g *ginContext) Request() *http.Request {
	if g.c == nil {
		return nil
	}
	return g.c.Request
}

// SetData implements Contexts.
func (g *ginContext) SetData(data interface{}) {
	g.data = data
}

func (c *ginContext) clone() *ginContext {
	ctx := osCtx.WithValue(c.ctx, "_", c.Value("_"))
	// copy all fields
	return &ginContext{
		ctx: ctx,
		c:   c.c,
		extraResponse: c.extraResponse,
		pageResponse: c.pageResponse,
		data:         c.data,
		requestID:    c.requestID,
		responseFile: c.responseFile,
		rawResponse:  c.rawResponse,
		cancelFn:     c.cancelFn,
		meErr:        c.meErr,
	}
}

 
 

// SubContext implements Contexts.
func (g *ginContext) SubContext(suffix string) Context {
	ctx := NewSubLogCtx(g.ctx, suffix)
	newCtx := g.clone()
	newCtx.ctx = ctx
	newCtx.requestID = GetLogID(ctx)
	return newCtx

}

// Value implements Contexts.
func (g *ginContext) Value(key any) any {
	panic("unimplemented")
}

// WithSpan implements Contexts.
func (g *ginContext) WithSpan() Context {
	panic("unimplemented")
}

// WithSpanID implements Contexts.
func (g *ginContext) WithSpanID(id string) Context {
	panic("unimplemented")
}

// WithSpanPrefix implements Contexts.
func (g *ginContext) WithSpanPrefix(prefix string) Context {
	panic("unimplemented")
}

// WithTimeoutCtx implements Contexts.
func (g *ginContext) WithTimeoutCtx(timeout time.Duration) Context {
	panic("unimplemented")
}

// WithValue implements Contexts.
func (g *ginContext) WithValue(key interface{}, value interface{}) {
	panic("unimplemented")
}

// NewGinContext creates a new context wrapper for gin.Context
func NewGinContext(c *gin.Context) Contexts {
	ctx := &ginContext{
		c:             c,
		extraResponse: make(map[string]interface{}),
		ctx:           c.Request.Context(),
		requestID:     c.GetString("requestID"), // You might want to generate this
		meErr:         &err{},
	}
	return ctx
}

var _ Contexts = (*ginContext)(nil)

func (g *ginContext) Log() Log {
	return &log{ctx: g}
}

// WithTimeout 设置超时间时间
func (c *ginContext) WithTimeout(timeout time.Duration) {
	c.ctx, c.cancelFn = osCtx.WithTimeout(c.ctx, timeout)
}

// AllMapper automatic data with struct private field
func (g *ginContext) AllMapper(action string, src, dst interface{}) errors.Error {
	if src == nil {
		return nil
	}
	if err := mapper.AllMapper(g.c, src, dst); err != nil {
		g.Log().ErrorJSON("action: %s, src: %s, err: %s", action, src, err.Error())
		return g.meErr.Errorf(code.MapperActionErrCode, action, err.Error())
	}

	return nil
}

func (g *ginContext) Response() http.ResponseWriter {
	return g.c.Writer
}

func (g *ginContext) JSONDecode(target interface{}) errors.Error {
	if err := decodeJSON(g.c.Request.Body, target); err != nil {
		return g.meErr.LegacyWrapCode(code.JSONDecodeErrCode, err)
	}

	if valid, ok := target.(ValidateI); ok {
		if err := valid.Validate(); err != nil {
			return err
		}
	}
	return nil

}

func (g *ginContext) Decode(target interface{}) errors.Error {
	pathParameters := g.PathParameters()
	urlParams := make(map[string][]string, 0)
	for k, v := range pathParameters {
		urlParams[k] = []string{v}
	}

	if err := autoDecode(g.c.Request, urlParams, target); err != nil {
		return g.meErr.LegacyWrapCode(code.JSONDecodeErrCode, err)
	}

	if valid, ok := target.(ValidateI); ok {

		if err := valid.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (g *ginContext) FromFile(name string) (multipart.File, *multipart.FileHeader, error) {
	return g.c.Request.FormFile(name)
}

func (g *ginContext) Query(name string) []string {
	values := g.c.Request.URL.Query()[name]
	if len(values) == 0 {
		// Also check form values
		if g.c.Request.Form != nil {
			values = g.c.Request.Form[name]
		}
	}
	return values
}

func (g *ginContext) PathParameter(name string) string {
	return g.c.Param(name)
}

func (g *ginContext) PathParameters() map[string]string {
	params := make(map[string]string)
	for _, param := range g.c.Params {
		params[param.Key] = param.Value
	}
	return params
}

func (g *ginContext) SetExtraResponse(key string, val interface{}) {
	g.extraResponse[key] = val
}

func (g *ginContext) SetPageResponse(val interface{}) {
	g.pageResponse = val
}

func (g *ginContext) SetResponseFile(fileName string, content *bytes.Buffer) {
	g.responseFile.fileName = fileName
	g.responseFile.content = content
	g.responseFile.exists = true
}

func (g *ginContext) GetResponseFile() (fileName string, content *bytes.Buffer, exists bool) {
	return g.responseFile.fileName, g.responseFile.content, g.responseFile.exists
}

func (g *ginContext) GetExtraResponse() map[string]interface{} {
	return g.extraResponse
}

func (g *ginContext) SetRawResponse(typ string, body []byte) {
	g.rawResponse.typ = typ
	g.rawResponse.body = body
	g.rawResponse.exists = true
}

func (g *ginContext) GetRawResponse() (typ string, body []byte, exists bool) {
	return g.rawResponse.typ, g.rawResponse.body, g.rawResponse.exists
}

func (g *ginContext) SelectedRoutePath() string {
	return g.c.FullPath()
}

// Additional gin-specific helper methods can be added here

func (g *ginContext) GetGinContext() *gin.Context {
	return g.c
}

func (g *ginContext) SetValue(key string, value interface{}) {
	g.c.Set(key, value)
}

func (g *ginContext) GetValue(key string) interface{} {
	value, exists := g.c.Get(key)
	if !exists {
		return nil
	}
	return value
}
