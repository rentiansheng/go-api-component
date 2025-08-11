package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/emicklei/go-restful/v3"
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
	container := h.initRoutes()

	// Create server with timeouts
	server := &http.Server{
		Addr:         ":" + h.s.Port,
		Handler:      container,
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

func (h *httpServer) initRoutes() *restful.Container {

	// Create container
	container := restful.NewContainer()

	if h.s.Cors.EnableCORS {
		s := h.s
		// Add CORS support
		cors := restful.CrossOriginResourceSharing{
			AllowedHeaders: s.Cors.AllowedHeaders, // []string{"Content-Type", "Accept", "Authorization"},
			AllowedMethods: s.Cors.AllowedMethods, // []string{"GET", "POST", "PUT", "DELETE"},
			AllowedDomains: s.Cors.AllowedDomains, //[]string{"*"},
			CookiesAllowed: s.Cors.CookiesAllowed, //true,
			Container:      container,
		}
		container.Filter(cors.Filter)
	}

	// Add logging filter
	container.Filter(logFilter)

	// Enable debugging
	restful.EnableTracing(true)
	restful.DefaultContainer.EnableContentEncoding(true)

	// Register routes
	registerRoutes(container, router.Get())

	return container
}

func registerRoutes(container *restful.Container, routers []router.Router) {
	for _, r := range routers {
		ws := r.Routes()
		container.Add(ws)
		for _, route := range ws.Routes() {
			log.Printf("Registered route: %s %s", route.Method, route.Path)
		}
	}
}

func logFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	start := time.Now()

	// Log request details
	log.Printf("[REQUEST] %s %s %s", req.Request.Method, req.Request.URL.Path, req.Request.RemoteAddr)

	chain.ProcessFilter(req, resp)

	duration := time.Since(start)
	log.Printf("[RESPONSE] %s %s %s %v", req.Request.Method, req.Request.URL.Path, req.Request.RemoteAddr, duration)
}
