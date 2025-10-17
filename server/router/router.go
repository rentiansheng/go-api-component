package router

import (
	"github.com/gin-gonic/gin"
)

var (
	routers []Router
)

type Router interface {
	RegisterGinRoutes(engine *gin.Engine)
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
