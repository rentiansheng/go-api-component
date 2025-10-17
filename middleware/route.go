package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	path2 "path"
)

type Web interface {
	Root(string)
	Get(string) Route
	Post(string) Route
	Put(string) Route
	Delete(string) Route
	Patch(string) Route
	Head(string) Route
	Options(string) Route

	Route(r Route)

	RegisterGinRoutes(engine *gin.Engine)
}

type Route interface {
	NoLogin() Route
	NeedLogin() Route
	Handler(h Handler) Route
	GetPath() string
	GetMethod() string
	GetHandler() Handler
	IsLoginRequired() bool
}

type ContentType string

const (
	ContentTypeJSON        ContentType = "application/json"
	ContentTypeXML         ContentType = "application/xml"
	ContentTypeZIP         ContentType = "application/zip"
	ContentTypeOctetStream ContentType = "application/octet-stream"
	ContentTypeProtoBuf    ContentType = "application/x-protobuf"
	ContentTypeMsgPack     ContentType = "application/x-msgpack"
	ContentTypeYaml        ContentType = "application/x-yaml"
	ContentTypeToml        ContentType = "application/toml"
)

func NewWeb(root string) Web {
	return &web{
		routers: make([]Route, 0),
		root:    root,
	}
}

type web struct {
	routers []Route
	root    string
}

func (w web) Get(path string) Route {
	o := route{}
	return o.Get(path)
}

func (w web) Post(path string) Route {
	o := route{}
	return o.Post(path)
}

func (w web) Put(path string) Route {
	o := route{}
	return o.Put(path)
}

func (w web) Delete(path string) Route {
	o := route{}
	return o.Delete(path)
}

func (w web) Patch(path string) Route {
	o := route{}
	return o.Patch(path)
}

func (w web) Head(path string) Route {
	o := route{}
	return o.Head(path)
}

func (w web) Options(path string) Route {
	o := route{}
	return o.Options(path)
}

func (w *web) Root(root string) {
	w.root = root
}

func (w *web) Route(r Route) {
	w.routers = append(w.routers, r)
type web struct {
	routers []Route
	root    string
}

func (w web) Get(path string) Route {
	r := &route{method: http.MethodGet, path: path}
	return r
}

func (w web) Post(path string) Route {
	r := &route{method: http.MethodPost, path: path}
	return r
}

func (w web) Put(path string) Route {
	r := &route{method: http.MethodPut, path: path}
	return r
}

func (w web) Delete(path string) Route {
	r := &route{method: http.MethodDelete, path: path}
	return r
}

func (w web) Patch(path string) Route {
	r := &route{method: http.MethodPatch, path: path}
	return r
}

func (w web) Head(path string) Route {
	r := &route{method: http.MethodHead, path: path}
	return r
}

func (w web) Options(path string) Route {
	r := &route{method: http.MethodOptions, path: path}
	return r
}

func (w *web) Root(root string) {
	w.root = root
}

func (w *web) Route(r Route) {
	w.routers = append(w.routers, r)
}

func (w *web) RegisterGinRoutes(engine *gin.Engine) {
	for _, r := range w.routers {
		fullPath := path2.Join(w.root, r.GetPath())
		handler := wrapGinHandler(r.GetHandler(), r.IsLoginRequired())
		
		switch r.GetMethod() {
		case http.MethodGet:
			engine.GET(fullPath, handler)
		case http.MethodPost:
			engine.POST(fullPath, handler)
		case http.MethodPut:
			engine.PUT(fullPath, handler)
		case http.MethodDelete:
			engine.DELETE(fullPath, handler)
		case http.MethodPatch:
			engine.PATCH(fullPath, handler)
		case http.MethodHead:
			engine.HEAD(fullPath, handler)
		case http.MethodOptions:
			engine.OPTIONS(fullPath, handler)
		}
	}
}

type route struct {
	noLogin     bool
	handler     Handler
	method      string
	path        string
	contentType ContentType
}

func (r *route) NoLogin() Route {
	r.noLogin = true
	return r
}

func (r *route) NeedLogin() Route {
	r.noLogin = false
	return r
}

func (r *route) Handler(h Handler) Route {
	r.handler = h
	return r
}

func (r *route) GetPath() string {
	return r.path
}

func (r *route) GetMethod() string {
	return r.method
}

func (r *route) GetHandler() Handler {
	return r.handler
}

func (r *route) IsLoginRequired() bool {
	return !r.noLogin
}

// wrapGinHandler converts our Handler to gin.HandlerFunc
func wrapGinHandler(handler Handler, needLogin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create context wrapper
		ctx := NewGinContext(c)
		
		// TODO: Add login check if needLogin is true
		// if needLogin {
		//     // Add authentication middleware logic here
		// }
		
		// Call the handler
		result, err := handler(ctx)
		
		// Handle the response
		if err != nil {
			// Handle error response
			c.JSON(err.Code(), gin.H{
				"code":    err.Code(),
				"message": err.Message(),
				"data":    nil,
			})
			return
		}
		
		// Handle success response
		c.JSON(200, gin.H{
			"code":    0,
			"message": "success",
			"data":    result,
		})
	}
}
