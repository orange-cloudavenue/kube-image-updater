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
	log "github.com/sirupsen/logrus"
)

const (
	timeout     = 10 * time.Second
	defaultPort = ":8080"
)

type (
	HTTPServer struct {
		Router *chi.Mux
		Config *http.Server
	}
	Option func(s *http.Server)
)

// NewHTTPServer returns a new HTTP router
// func New(path, port string, tlsC *tls.Config) (s HTTPServer) {
func New(opts ...Option) *HTTPServer {
	s := &HTTPServer{}
	s.Router = chi.NewRouter()
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Timeout(timeout))

	// Default server configuration
	s.Config = &http.Server{
		Addr:        defaultPort,
		Handler:     s.Router,
		ReadTimeout: timeout,
	}
	for _, opt := range opts {
		opt(s.Config)
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

// Add Get routes to the HTTP server
func (s HTTPServer) AddGetRoutes(path string, handler http.Handler) {
	s.Router.Mount(path, handler)
}

// Add Post routes to the HTTP server
func (s HTTPServer) AddPostRoutes(path string, handler http.Handler) {
	s.Router.Mount(path, handler)
}

// ServeHTTP implements the http.Handler interface
func (s HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// ListenAndServe starts the HTTP server
func (s HTTPServer) Start(ctx context.Context, wg *sync.WaitGroup) (err error) {
	wg.Add(1)
	defer wg.Done()

	switch s.Config.TLSConfig {
	case nil:
		// Start the HTTP server
		go func() {
			log.Infof("Starting server on %s", s.Config.Addr)
			// log.Printf("Starting server on %s", s.Config.Addr)
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
