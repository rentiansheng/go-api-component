package router

import (
	"github.com/emicklei/go-restful/v3"
)

var (
	routers []Router
)

type Router interface {
	Routes() *restful.WebService
}

func Register(r Router) {
	if r == nil {
		return
	}
	routers = append(routers, r)
}

func Get() []Router {
	return routers
}
