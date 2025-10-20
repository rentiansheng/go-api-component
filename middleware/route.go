package middleware

import (
	"net/http"
	path2 "path"

	"github.com/gin-gonic/gin"
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
	r := &route{}
	return r.Get(path)
}

func (w web) Post(path string) Route {
	r := &route{}
	return r.Post(path)
}

func (w web) Put(path string) Route {
	r := &route{}
	return r.Put(path)

}

func (w web) Delete(path string) Route {
	r := &route{}
	return r.Delete(path)
}

func (w web) Patch(path string) Route {
	r := &route{}
	return r.Patch(path)
}

func (w web) Head(path string) Route {
	r := &route{}
	return r.Head(path)
}

func (w web) Options(path string) Route {
	r := &route{}
	return r.Options(path)
}

func (w *web) Root(root string) {
	w.root = root
}

func (w *web) Route(r Route) {
	w.routers = append(w.routers, r)
}

func (w *web) RegisterGinRoutes(engine *gin.Engine) {
	for _, r := range w.routers {
		o := DefaultOption()
		if !r.IsLoginRequired() {
			o.WithNoLogin()
		}

		fullPath := path2.Join(w.root, r.GetPath())
		handler := wrapperOptions(r.GetHandler(), o)

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

func (r route) Get(path string) Route {
	r.method = http.MethodGet
	return r.Path(path)
}

func (r route) Post(path string) Route {
	r.method = http.MethodPost
	return r.Path(path)
}

func (r route) Put(path string) Route {
	r.method = http.MethodPut
	return r.Path(path)
}

func (r route) Delete(path string) Route {
	r.method = http.MethodDelete
	return r.Path(path)
}

func (r route) Patch(path string) Route {
	r.method = http.MethodPatch
	return r.Path(path)
}

func (r route) Head(path string) Route {
	r.method = http.MethodHead
	return r.Path(path)
}

func (r route) Options(path string) Route {
	r.method = http.MethodOptions
	return r.Path(path)
}

func (r *route) Path(path string) Route {
	r.path = path
	return r
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
