package httpserver

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/orange-cloudavenue/kube-image-updater/internal/health"
)

var _ InterfaceServer = &app{}

type (

	// InterfaceSerever is an interface to manage the server
	// Run: Start the server
	// Add: Add a new server
	InterfaceServer interface {
		// Start all server listed in the app struct
		Run() (err error)
		// Add a new server to the app without endpoint.
		// return a server where you can add endpoint GET, POST, PUT, DELETE...
		Add(name string, opts ...Option) (s *server, err error)
	}

	app struct {
		list map[string]*server
		ctx  context.Context
		wg   *sync.WaitGroup
	}

	server struct {
		// Set http parameters
		http *http.Server
		// Configure endpoint http
		Config *chi.Mux
	}

	// OptionsHTTP is a function to set the http server
	Option       func(s *server)
	OptionServer func(a *app)

	CancelFunc func()
)

var (
	defaultPortHealth  string = ":9081"
	defaultPathHealth  string = "/healthz"
	defaultPortMetrics string = ":9080"
	defaultPathMetrics string = "/metrics"
	defaultAddr        string = ":8080"
	timeoutR                  = 5 * time.Second
)

func init() {
	flag.StringVar(&defaultPortHealth, "health-port", defaultPortHealth, "Health server port. ex: :9081")
	flag.StringVar(&defaultPathHealth, "health-path", defaultPathHealth, "Health server path. ex: /healthz")
	flag.StringVar(&defaultPortMetrics, "metrics-port", defaultPortMetrics, "Metrics server port. ex: :9080")
	flag.StringVar(&defaultPathMetrics, "metrics-path", defaultPathMetrics, "Metrics server path. ex: /metrics")
}

// Function to initialize application, return app struct and a func waitgroup.
// The app contains a map of server.
// By default, the app contains a health and metrics server.
func Init(ctx context.Context, opts ...OptionServer) (InterfaceServer, CancelFunc) {
	a := &app{
		list: make(map[string]*server),
		ctx:  ctx,
		wg:   &sync.WaitGroup{},
	}

	a.list["health"] = a.createHealth()
	a.list["metrics"] = a.createMetrics()

	// create a new server for health
	for _, opt := range opts {
		opt(a)
	}

	return a, func() {
		log.Info("Waiting for all server to shutdown gracefully...")
		a.wg.Wait()
		log.Info("All Server on has been shutdown: bye...")
	}
}

// Function to disable the health server
func DisableHealth() OptionServer {
	return func(a *app) {
		delete(a.list, "health")
	}
}

// Function to disable the metrics server
func DisableMetrics() OptionServer {
	return func(a *app) {
		delete(a.list, "metrics")
	}
}

// Function to create a new server for health
func (a *app) createHealth() *server {
	s := a.new(WithAddr(defaultPortHealth))
	s.Config.Get(defaultPathHealth, health.Handler().ServeHTTP)
	return s
}

// Function to create a new server for metrics
func (a *app) createMetrics() *server {
	s := a.new(WithAddr(defaultPortMetrics))
	s.Config.Get(defaultPathMetrics, promhttp.Handler().ServeHTTP)
	return s
}

// Function return a server
func (a *app) new(opts ...Option) *server {
	// create a new router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// create a new server with default parameters
	s := &server{
		http: &http.Server{
			Addr:        defaultAddr,
			ReadTimeout: timeoutR,
		},
		Config: r,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Function to set the address of the server
// addr must be like "ipv4:1-65535"
// addr must be an IPV4 format and with a port number between 1 and 65535
// if addr is not correct, the default local address is set
func WithAddr(addr string) Option {
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		return func(s *server) {
			s.http.Addr = defaultAddr
		}
	}
	return func(s *server) {
		s.http.Addr = addr
	}
}

// Function to set the TLS configuration
func WithTLS(tlsC *tls.Config) Option {
	return func(s *server) {
		s.http.TLSConfig = tlsC
	}
}

// Function to start the server
func (a *app) Run() (err error) {
	for name, s := range a.list {
		if s.http.TLSConfig != nil {
			log.Infof("Starting server %s on %s with TLS", name, s.http.Addr)
		} else {
			log.Infof("Starting server %s on %s", name, s.http.Addr)
		}
		if err = a.start(s); err != nil {
			return err
		}
	}
	return nil
}

// ListenAndServe starts the HTTP server
func (a *app) start(s *server) (err error) {
	a.wg.Add(1)
	defer a.wg.Done()

	switch s.http.TLSConfig {
	case nil:
		// Start the HTTP server
		go func() {
			s.http.Handler = s.Config
			if err = s.http.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				return
			}
		}()

	default:
		// Start the HTTPS server
		go func() {
			s.http.Handler = s.Config
			if err = s.http.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) {
				return
			}
		}()
	}

	// Kill the server if there is an error or stop signal
	go func() {
		for {
			<-a.ctx.Done()
			a.wg.Add(1)
			defer a.wg.Done()
			ctxTimeout, cancel := context.WithTimeout(a.ctx, 5*time.Second)
			log.Infof("Shutting down server on %s", s.http.Addr)
			cancel()
			if err = s.http.Shutdown(ctxTimeout); err != nil {
				log.Errorf("Failed to shutdown HTTP server on %s: %v", s.http.Addr, err)
			}
			return
		}
	}()
	return nil
}

// Add a new server to the app without endpoint
// return a server where you can add endpoint GET, POST, PUT, DELETE...
func (a *app) Add(name string, opts ...Option) (s *server, err error) {
	s = a.new(opts...)
	if a.checkIfPortIsAlreadyUsed(s) {
		return nil, fmt.Errorf("port %s is already used", s.http.Addr)
	}
	a.list[name] = s
	return
}

// Function to check if the port is already used
func (a *app) checkIfPortIsAlreadyUsed(s *server) bool {
	for _, v := range a.list {
		if v.http.Addr == s.http.Addr {
			return true
		}
	}
	return false
}
