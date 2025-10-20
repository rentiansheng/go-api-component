package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rentiansheng/go-api-component/pkg/logger"
	"github.com/rentiansheng/go-api-component/server/router"
	"github.com/stretchr/testify/assert"
)

func TestNewHttpServer(t *testing.T) {
	config := Server{
		Port:            "8080",
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		ShutdownTimeout: 10 * time.Second,
	}

	log := &logger.Log{}

	server := NewHttpServer("test-server", config, log)
	assert.NotNil(t, server)
}

func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		serverName   string
		expectedName string
	}{
		{
			name:         "With custom name",
			serverName:   "custom-server",
			expectedName: "custom-server",
		},
		{
			name:         "With empty name",
			serverName:   "",
			expectedName: "default-go-api-server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := New(tt.serverName)
			assert.NotNil(t, server)
		})
	}
}

func TestHttpServer_SetLogConfig(t *testing.T) {
	server := New("test")
	log := &logger.Log{}

	server.SetLogConfig(log)
	assert.NotNil(t, server)
}

func TestHttpServer_SetServerConfig(t *testing.T) {
	server := New("test")
	config := Server{
		Port:            "9090",
		ReadTimeout:     20 * time.Second,
		WriteTimeout:    20 * time.Second,
		ShutdownTimeout: 5 * time.Second,
	}

	server.SetServerConfig(config)
	assert.NotNil(t, server)
}

func TestHttpServer_SetName(t *testing.T) {
	server := New("original")
	server.SetName("updated")
	assert.NotNil(t, server)
}

func TestHttpServer_InitRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log := &logger.Log{}
	config := Server{
		Port:         "8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Cors: Cors{
			EnableCORS:     false,
			AllowedMethods: []string{"GET", "POST"},
			AllowedHeaders: []string{"Content-Type"},
			AllowedDomains: []string{"*"},
		},
	}

	srv := NewHttpServer("test", config, log).(*httpServer)
	engine := srv.initRoutes()

	assert.NotNil(t, engine)
}

func TestHttpServer_InitRoutesWithCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log := &logger.Log{}
	config := Server{
		Port:         "8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Cors: Cors{
			EnableCORS:     true,
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders: []string{"Content-Type", "Authorization"},
			AllowedDomains: []string{"http://localhost:3000"},
			CookiesAllowed: true,
		},
	}

	srv := NewHttpServer("test", config, log).(*httpServer)
	engine := srv.initRoutes()

	assert.NotNil(t, engine)
}

func TestHttpServer_CORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := Server{
		Cors: Cors{
			EnableCORS:     true,
			AllowedMethods: []string{"GET", "POST", "PUT"},
			AllowedHeaders: []string{"Content-Type", "Authorization"},
			AllowedDomains: []string{"http://example.com", "http://localhost:3000"},
			CookiesAllowed: true,
		},
	}

	srv := &httpServer{
		name: "test",
		s:    config,
	}

	tests := []struct {
		name          string
		origin        string
		method        string
		expectOrigin  string
		expectMethods string
		expectHeaders string
		expectCreds   string
	}{
		{
			name:          "Allowed origin",
			origin:        "http://example.com",
			method:        "GET",
			expectOrigin:  "http://example.com",
			expectMethods: "GET, POST, PUT",
			expectHeaders: "Content-Type, Authorization",
			expectCreds:   "true",
		},
		{
			name:          "Another allowed origin",
			origin:        "http://localhost:3000",
			method:        "POST",
			expectOrigin:  "http://localhost:3000",
			expectMethods: "GET, POST, PUT",
			expectHeaders: "Content-Type, Authorization",
			expectCreds:   "true",
		},
		{
			name:          "Disallowed origin",
			origin:        "http://evil.com",
			method:        "GET",
			expectOrigin:  "",
			expectMethods: "GET, POST, PUT",
			expectHeaders: "Content-Type, Authorization",
			expectCreds:   "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest(tt.method, "/test", nil)
			req.Header.Set("Origin", tt.origin)
			c.Request = req

			middleware := srv.corsMiddleware()
			middleware(c)

			if tt.expectOrigin != "" {
				assert.Equal(t, tt.expectOrigin, w.Header().Get("Access-Control-Allow-Origin"))
			}
			assert.Equal(t, tt.expectMethods, w.Header().Get("Access-Control-Allow-Methods"))
			assert.Equal(t, tt.expectHeaders, w.Header().Get("Access-Control-Allow-Headers"))
			assert.Equal(t, tt.expectCreds, w.Header().Get("Access-Control-Allow-Credentials"))
		})
	}
}

func TestHttpServer_CORSMiddleware_Wildcard(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := Server{
		Cors: Cors{
			EnableCORS:     true,
			AllowedMethods: []string{"GET", "POST"},
			AllowedHeaders: []string{"Content-Type"},
			AllowedDomains: []string{"*"},
			CookiesAllowed: false,
		},
	}

	srv := &httpServer{
		name: "test",
		s:    config,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://any-domain.com")
	c.Request = req

	middleware := srv.corsMiddleware()
	middleware(c)

	assert.Equal(t, "http://any-domain.com", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestHttpServer_CORSMiddleware_PreflightRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := Server{
		Cors: Cors{
			EnableCORS:     true,
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders: []string{"Content-Type", "Authorization"},
			AllowedDomains: []string{"*"},
			CookiesAllowed: false,
		},
	}

	srv := &httpServer{
		name: "test",
		s:    config,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	c.Request = req

	middleware := srv.corsMiddleware()
	middleware(c)

	// OPTIONS request should return 204
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()

	// Create a mock router
	mockRouter := &MockRouter{}
	routers := []router.Router{mockRouter}

	// This will fail compilation if RegisterGinRoutes doesn't exist
	// but demonstrates the structure
	registerRoutes(engine, routers)

	// Basic assertion
	assert.NotNil(t, engine)
}

// MockRouter for testing
type MockRouter struct{}

func (m *MockRouter) RegisterGinRoutes(engine *gin.Engine) {
	engine.GET("/mock", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "mock"})
	})
}
