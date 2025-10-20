package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/rentiansheng/go-api-component/middleware/context"
	. "github.com/rentiansheng/go-api-component/middleware/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewWeb(t *testing.T) {
	web := NewWeb("/api/v1")
	assert.NotNil(t, web)
}

func TestWeb_Get(t *testing.T) {
	web := NewWeb("/api/v1")
	route := web.Get("/users")

	assert.NotNil(t, route)
	assert.Equal(t, "/users", route.GetPath())
	assert.Equal(t, http.MethodGet, route.GetMethod())
}

func TestWeb_Post(t *testing.T) {
	web := NewWeb("/api/v1")
	route := web.Post("/users")

	assert.NotNil(t, route)
	assert.Equal(t, "/users", route.GetPath())
	assert.Equal(t, http.MethodPost, route.GetMethod())
}

func TestWeb_Put(t *testing.T) {
	web := NewWeb("/api/v1")
	route := web.Put("/users/:id")

	assert.NotNil(t, route)
	assert.Equal(t, "/users/:id", route.GetPath())
	assert.Equal(t, http.MethodPut, route.GetMethod())
}

func TestWeb_Delete(t *testing.T) {
	web := NewWeb("/api/v1")
	route := web.Delete("/users/:id")

	assert.NotNil(t, route)
	assert.Equal(t, "/users/:id", route.GetPath())
	assert.Equal(t, http.MethodDelete, route.GetMethod())
}

func TestWeb_Patch(t *testing.T) {
	web := NewWeb("/api/v1")
	route := web.Patch("/users/:id")

	assert.NotNil(t, route)
	assert.Equal(t, "/users/:id", route.GetPath())
	assert.Equal(t, http.MethodPatch, route.GetMethod())
}

func TestWeb_Head(t *testing.T) {
	web := NewWeb("/api/v1")
	route := web.Head("/users")

	assert.NotNil(t, route)
	assert.Equal(t, "/users", route.GetPath())
	assert.Equal(t, http.MethodHead, route.GetMethod())
}

func TestWeb_Options(t *testing.T) {
	web := NewWeb("/api/v1")
	route := web.Options("/users")

	assert.NotNil(t, route)
	assert.Equal(t, "/users", route.GetPath())
	assert.Equal(t, http.MethodOptions, route.GetMethod())
}

func TestRoute_NoLogin(t *testing.T) {
	web := NewWeb("/api/v1")
	route := web.Get("/public").NoLogin()

	assert.NotNil(t, route)
	assert.False(t, route.IsLoginRequired())
}

func TestRoute_NeedLogin(t *testing.T) {
	web := NewWeb("/api/v1")
	route := web.Get("/private").NeedLogin()

	assert.NotNil(t, route)
	assert.True(t, route.IsLoginRequired())
}

func TestRoute_Handler(t *testing.T) {
	web := NewWeb("/api/v1")

	handler := func(ctx Contexts) Error {
		ctx.SetData(map[string]string{"message": "success"})
		return nil
	}

	route := web.Get("/test").Handler(handler)

	assert.NotNil(t, route)
	assert.NotNil(t, route.GetHandler())
}

func TestWeb_RegisterGinRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	web := NewWeb("/api/v1")

	// Add some routes
	r := web.Get("/users").NoLogin().Handler(func(ctx Contexts) Error {
		ctx.SetData([]string{"user1", "user2"})
		return nil
	})
	web.Route(r)

	r = web.Post("/users").NeedLogin().Handler(func(ctx Contexts) Error {
		ctx.SetData(map[string]string{"id": "123"})
		return nil
	})
	web.Route(r)

	r = web.Get("/users/:id").NoLogin().Handler(func(ctx Contexts) Error {
		id := ctx.PathParameter("id")
		ctx.SetData(map[string]string{"id": id})
		return nil
	})
	web.Route(r)

	// Register routes to gin engine
	engine := gin.New()
	web.RegisterGinRoutes(engine)

	// Test GET /api/v1/users
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test POST /api/v1/users
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/v1/users", nil)
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test GET /api/v1/users/123
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/v1/users/123", nil)
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWeb_Root(t *testing.T) {
	web := NewWeb("/api/v1").(*web)
	assert.Equal(t, "/api/v1", web.root)

	web.Root("/api/v2")
	assert.Equal(t, "/api/v2", web.root)
}

func TestWeb_Route(t *testing.T) {
	web := NewWeb("/api/v1").(*web)

	route := web.Get("/test").Handler(func(ctx Contexts) Error {
		return nil
	})

	initialLen := len(web.routers)
	web.Route(route)
	assert.Equal(t, initialLen+1, len(web.routers))
}

func TestWrapGinHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := func(ctx Contexts) Error {
		ctx.SetData(map[string]string{"status": "ok"})
		return nil
	}

	ginHandler := wrapperOptions(handler, DefaultOption())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/test", nil)
	c.Request = req

	ginHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestWrapGinHandler_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := func(ctx Contexts) Error {
		return NewError(400, "Bad request")
	}

	ginHandler := wrapperOptions(handler, DefaultOption())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/test", nil)
	c.Request = req

	ginHandler(c)

	assert.Contains(t, w.Body.String(), "Bad request")
}

func TestRoute_ChainedCalls(t *testing.T) {
	web := NewWeb("/api/v1")

	route := web.Get("/test").
		NoLogin().
		Handler(func(ctx Contexts) Error {
			ctx.SetData("test")
			return nil
		})

	assert.NotNil(t, route)
	assert.Equal(t, "/test", route.GetPath())
	assert.Equal(t, http.MethodGet, route.GetMethod())
	assert.False(t, route.IsLoginRequired())
	assert.NotNil(t, route.GetHandler())
}

func TestRoute_AllMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)

	web := NewWeb("/api/v1")
	engine := gin.New()

	handler := func(ctx Contexts) Error {
		ctx.SetData(map[string]string{"method": ctx.Request().Method})
		return nil
	}

	// Register all HTTP methods
	r := web.Get("/resource").NoLogin().Handler(handler)
	web.Route(r)
	r = web.Post("/resource").NoLogin().Handler(handler)
	web.Route(r)
	r = web.Put("/resource").NoLogin().Handler(handler)
	web.Route(r)
	r = web.Delete("/resource").NoLogin().Handler(handler)
	web.Route(r)
	r = web.Patch("/resource").NoLogin().Handler(handler)
	web.Route(r)
	r = web.Head("/resource").NoLogin().Handler(handler)
	web.Route(r)
	r = web.Options("/resource").NoLogin().Handler(handler)
	web.Route(r)

	web.RegisterGinRoutes(engine)

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}

	for _, method := range methods {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(method, "/api/v1/resource", nil)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Method: "+method)
	}
}

func TestContentTypeConstants(t *testing.T) {
	assert.Equal(t, ContentType("application/json"), ContentTypeJSON)
	assert.Equal(t, ContentType("application/xml"), ContentTypeXML)
	assert.Equal(t, ContentType("application/zip"), ContentTypeZIP)
	assert.Equal(t, ContentType("application/octet-stream"), ContentTypeOctetStream)
	assert.Equal(t, ContentType("application/x-protobuf"), ContentTypeProtoBuf)
	assert.Equal(t, ContentType("application/x-msgpack"), ContentTypeMsgPack)
	assert.Equal(t, ContentType("application/x-yaml"), ContentTypeYaml)
	assert.Equal(t, ContentType("application/toml"), ContentTypeToml)
}
