package httpserver

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/orange-cloudavenue/kube-image-updater/internal/health"
)

const (
	timeout                   = 10 * time.Second
	timeoutR                  = 10 * time.Second
	timeoutW                  = 10 * time.Second
	defaultPort               = ":8080"
	defaultPortMetrics string = ":9080"
	defaultPortHealth  string = ":9081"
	defaultPathMetrics string = "/metrics"
	defaultPathHealth  string = "/healthz"
)

var wg *sync.WaitGroup

type (
	HTTPServer struct {
		Router *chi.Mux
		Config *http.Server
	}
	Option        func(s *http.Server)
	OptionMetrics func(port, path string)
)

// Func Init() initialize the waitgroup
// return a func to wait all server to shutdown gracefully
func Init() (waitStop func()) {
	wg = &sync.WaitGroup{}
	return WaitStop
}

// NewHTTPServer returns a new HTTP router
// func New(path, port string, tlsC *tls.Config) (s HTTPServer) {
func New(opts ...Option) *HTTPServer {
	s := &HTTPServer{}
	s.Router = chi.NewRouter()
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Timeout(timeout))

	// Default server configuration
	s.Config = &http.Server{
		Addr:         defaultPort,
		Handler:      s.Router,
		ReadTimeout:  timeoutR,
		WriteTimeout: timeoutW,
	}
	for _, opt := range opts {
		opt(s.Config)
	}
	// check if waitgroup exist
	// if not, create a new one
	if wg == nil {
		wg = &sync.WaitGroup{}
	}
	return s
}

// WithTLSConfig sets the TLS configuration for the HTTP server
// Add an option to set the TLS configuration for the HTTP server
// The WithTLSConfig function takes a *tls.Config as an argument and returns an Option
// The Option type is a function that takes a *http.Server as an argument
//
// ex: New(httpserver.WithTLSConfig(tlsC))
// ex: New(httpserver.WithTLSConfig(tlsC), httpserver.WithAddr(":8443"))
func WithTLSConfig(tlsC *tls.Config) Option {
	return func(s *http.Server) {
		s.TLSConfig = tlsC
	}
}

// WithAddr sets the address for the HTTP server
// Add an option to set the address for the HTTP server
// The WithAddr function takes a string as an argument and returns an Option
// The Option type is a function that takes a *http.Server as an argument
//
// ex: New(httpserver.WithAddr(":8443"))
// ex: New(httpserver.WithTLSConfig(tlsC), httpserver.WithAddr(":8443"))
func WithAddr(addr string) Option {
	return func(s *http.Server) {
		s.Addr = addr
	}
}

// WithTimeout sets the timeout for the HTTP server
// Add an option to set the timeout for the HTTP server
// The WithTimeout function takes a time.Duration as an argument and returns an Option
// The Option type is a function that takes a *http.Server as an argument
//
// ex: New(httpserver.WithTimeout(10*time.Second))
// ex: New(httpserver.WithTLSConfig(tlsC), httpserver.WithAddr(":8443"), httpserver.WithTimeout(10*time.Second))
func WithReadTimeout(timeout time.Duration) Option {
	return func(s *http.Server) {
		s.ReadTimeout = timeout
	}
}

// WithWriteTimeout sets the write timeout for the HTTP server
// Add an option to set the write timeout for the HTTP server
// The WithWriteTimeout function takes a time.Duration as an argument and returns an Option
// The Option type is a function that takes a *http.Server as an argument
//
// ex: New(httpserver.WithWriteTimeout(10*time.Second))
// ex: New(httpserver.WithTLSConfig(tlsC), httpserver.WithAddr(":8443"), httpserver.WithWriteTimeout(10*time.Second))
func WithWriteTimeout(timeout time.Duration) Option {
	return func(s *http.Server) {
		s.WriteTimeout = timeout
	}
}

// WithHandler sets the handler for the HTTP server
// Add an option to set the handler for the HTTP server
// The WithHandler function takes a http.Handler as an argument and returns an Option
// The Option type is a function that takes a *http.Server as an argument
//
// ex: New(httpserver.WithHandler(handler))
// ex: New(httpserver.WithTLSConfig(tlsC), httpserver.WithAddr(":8443"), httpserver.WithHandler(handler))
func WithHandler(handler http.Handler) Option {
	return func(s *http.Server) {
		s.Handler = handler
	}
}

// Add Get routes to the HTTP server
func (s *HTTPServer) AddGetRoutes(path string, handler http.Handler) {
	s.Router.Mount(path, handler)
}

// Add Post routes to the HTTP server
func (s *HTTPServer) AddPostRoutes(path string, handler http.Handler) {
	s.Router.Mount(path, handler)
}

// ListenAndServe starts the HTTP server
func (s *HTTPServer) Start(ctx context.Context) (err error) {
	wg.Add(1)
	defer wg.Done()

	switch s.Config.TLSConfig {
	case nil:
		// Start the HTTP server
		go func() {
			log.Infof("Starting server on %s", s.Config.Addr)
			if err = s.Config.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				return
			}
		}()

	default:
		// Start the HTTPS server
		go func() {
			log.Infof("Starting TLS server on %s", s.Config.Addr)
			if err = s.Config.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) {
				return
			}
		}()
	}

	// Kill the server if there is an error or stop signal
	go func() {
		for {
			<-ctx.Done()
			wg.Add(1)
			defer wg.Done()
			ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
			log.Infof("Shutting down server on %s", s.Config.Addr)
			cancel()
			if err = s.Config.Shutdown(ctxTimeout); err != nil {
				log.Errorf("Failed to shutdown HTTP server on %s: %v", s.Config.Addr, err)
			}
			return
		}
	}()
	return nil
}

// StartMetrics starts the HTTP server for Metrics
// With default port (:9080) and path /metrics
func StartMetrics(ctx context.Context, opts ...Option) (err error) {
	s := New(WithAddr(defaultPortMetrics))
	s.AddGetRoutes(defaultPathMetrics, promhttp.Handler())
	return s.Start(ctx)
}

// StartHealth starts the HTTP server for Health
// With default port (:9081) and path /healthz
func StartHealth(ctx context.Context, opts ...Option) (err error) {
	s := New(WithAddr(defaultPortHealth))
	s.AddGetRoutes(defaultPathHealth, health.Handler())
	return s.Start(ctx)
}

// func use to wait ALL HTTP server to shutdown gracefully
func WaitStop() {
	log.Info("Waiting for all server to shutdown gracefully...")
	wg.Wait()
	log.Info("All Server on has been shutdown: bye...")
}
