package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rentiansheng/go-api-component/pkg/logger"
	"github.com/rentiansheng/go-api-component/server/router"
)

type Server struct {
	Port            string        `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	Cors            Cors          `mapstructure:"cors"`
}

type Cors struct {
	AllowedHeaders []string `mapstructure:"allowed_headers"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedDomains []string `mapstructure:"allowed_domains"`
	CookiesAllowed bool     `mapstructure:"cookies_allowed"`
	EnableCORS     bool     `mapstructure:"enable_cors"`
}

type httpServer struct {
	name string
	s    Server
	l    *logger.Log
}

type HTTPServer interface {
	SetLogConfig(l *logger.Log)
	SetServerConfig(s Server)
	SetName(name string)
	Run() error
}

func NewHttpServer(name string, s Server, l *logger.Log) HTTPServer {
	return &httpServer{
		name: name,
		s:    s,
		l:    l,
	}
}

func New(name string) HTTPServer {
	if name == "" {
		name = "default-go-api-server"
	}
	return &httpServer{
		name: name,
		s: Server{
			Port:            "8080",
			ReadTimeout:     10 * time.Second,
			WriteTimeout:    10 * time.Second,
			ShutdownTimeout: 5 * time.Second,
		},
		l: nil,
	}
}

func (h *httpServer) SetLogConfig(l *logger.Log) {
	h.l = l
}

func (h *httpServer) SetServerConfig(s Server) {
	h.s = s
}

func (h *httpServer) SetName(name string) {
	h.name = name
}

// Run initializes and starts the HTTP server
func (h *httpServer) Run() error {

	// Configure logging
	if h.l != nil {
		logger.Config(h.name, *h.l)
	}
	
	router := h.initRoutes()

	// Create server with timeouts
	server := &http.Server{
		Addr:         ":" + h.s.Port,
		Handler:      router,
		ReadTimeout:  h.s.ReadTimeout,
		WriteTimeout: h.s.WriteTimeout,
	}

	stop := make(chan os.Signal, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println("failed to start server:", err)
		}
		stop <- syscall.SIGABRT
		return
	}()

	// 捕获退出信号
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop // 阻塞直到退出

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return server.Shutdown(ctx)
}

func (h *httpServer) initRoutes() *gin.Engine {
	// Set gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create gin engine
	engine := gin.New()

	// Add middleware
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	if h.s.Cors.EnableCORS {
		// Add CORS middleware
		engine.Use(h.corsMiddleware())
	}

	// Register routes
	registerRoutes(engine, router.Get())

	return engine
}

func (h *httpServer) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := h.s.Cors
		
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// Check if origin is allowed
			allowed := false
			for _, domain := range s.AllowedDomains {
				if domain == "*" || domain == origin {
					allowed = true
					break
				}
			}
			
			if allowed {
				c.Header("Access-Control-Allow-Origin", origin)
			}
		}
		
		// Set other CORS headers
		if len(s.AllowedMethods) > 0 {
			methods := ""
			for i, method := range s.AllowedMethods {
				if i > 0 {
					methods += ", "
				}
				methods += method
			}
			c.Header("Access-Control-Allow-Methods", methods)
		}
		
		if len(s.AllowedHeaders) > 0 {
			headers := ""
			for i, header := range s.AllowedHeaders {
				if i > 0 {
					headers += ", "
				}
				headers += header
			}
			c.Header("Access-Control-Allow-Headers", headers)
		}
		
		if s.CookiesAllowed {
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		
		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

func registerRoutes(engine *gin.Engine, routers []router.Router) {
	for _, r := range routers {
		r.RegisterGinRoutes(engine)
		log.Printf("Registered router: %T", r)
	}
}
