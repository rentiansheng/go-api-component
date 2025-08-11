package middleware

import (
	"github.com/emicklei/go-restful/v3"
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

	Routes() *restful.WebService
}

type Route interface {
	NoLogin() Route
	NeedLogin() Route
	Handler(h Handler) Route
	RestfulRoute(root string) *restful.RouteBuilder
}

type ContentType string

const (
	ContentTypeJSON        ContentType = restful.MIME_JSON
	ContentTypeXML         ContentType = restful.MIME_XML
	ContentTypeZIP         ContentType = restful.MIME_ZIP
	ContentTypeOctetStream ContentType = restful.MIME_OCTET
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
}

func (w *web) Routes() *restful.WebService {
	web := &restful.WebService{}
	web.Path(w.root)
	for _, r := range w.routers {
		web.Route(r.RestfulRoute(w.root))
	}

	return web
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

// 默认是需要登录的
func (r route) NoLogin() Route {
	r.noLogin = true
	return r
}

func (r route) NeedLogin() Route {
	r.noLogin = false
	return r
}

func (r route) Path(path string) Route {
	r.path = path
	return r
}

func (r route) Handler(h Handler) Route {
	r.handler = h
	return r
}

func (r route) Produces(contentType ContentType) Route {
	r.contentType = contentType
	return r
}

func (r route) RestfulRoute(root string) *restful.RouteBuilder {

	o := DefaultOption()
	if r.noLogin {
		o = o.WithNoLogin()
	}
	path := r.path
	if len(root) != 0 {
		path = path2.Join(root, path)
	}
	rb := &restful.RouteBuilder{}
	rb = rb.Method(r.method).Path(path).To(wrapperOptions(r.handler, o))
	if len(r.contentType) != 0 {
		rb = rb.Produces(string(r.contentType))
	} else {
		rb = rb.Produces(restful.MIME_JSON)
	}
	return rb
}
