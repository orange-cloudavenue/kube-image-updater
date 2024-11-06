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

	logger "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/orange-cloudavenue/kube-image-updater/internal/log"
	"github.com/orange-cloudavenue/kube-image-updater/internal/metrics"
	"github.com/orange-cloudavenue/kube-image-updater/internal/models"
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
	Option func(s *server)
	// OptionsServer is a function to disable some server
	OptionServer func(a *app)

	CancelFunc func()

	// Function to check the health of the application
	HealthzFunc func() (bool, error)
)

var (
	HealthzPort int    = 0  // expose var to be able to operator
	HealthzPath string = "" // expose var to be able to operator
	metricsPort int    = 0
	metricsPath string = ""

	timeoutR                       = 5 * time.Second
	DefaultFuncHealthz HealthzFunc = func() (bool, error) {
		_, err := net.DialTimeout("tcp", models.HealthzDefaultAddr, timeoutR)
		if err != nil {
			return false, err
		}
		return true, nil
	}
)

func init() {
	// * Healthz
	flag.Bool(models.HealthzFlagName, false, "Enable the healthz server.")
	flag.IntVar(&HealthzPort, models.HealthzPortFlagName, int(models.HealthzDefaultPort), "Healthz server port.")
	flag.StringVar(&HealthzPath, models.HealthzPathFlagName, models.HealthzDefaultPath, "Healthz server path.")

	// * Metrics
	flag.Bool(models.MetricsFlagName, false, "Enable the metrics server.")
	flag.IntVar(&metricsPort, models.MetricsPortFlagName, int(models.MetricsDefaultPort), "Metrics server port.")
	flag.StringVar(&metricsPath, models.MetricsPathFlagName, models.MetricsDefaultPath, "Metrics server path.")
}

// Function to initialize application, return app struct and a func waitgroup.
// The app contains a map of server.
// By default, the app contains a healthz and metrics server.
func Init(ctx context.Context, opts ...OptionServer) (InterfaceServer, CancelFunc) {
	a := &app{
		list: make(map[string]*server),
		ctx:  ctx,
		wg:   &sync.WaitGroup{},
	}

	if flag.Lookup(models.HealthzFlagName).Value.String() == "true" {
		a.list["health"] = a.createHealth()
		WithCustomHandlerForHealth(DefaultFuncHealthz)(a)
	}

	if flag.Lookup(models.MetricsFlagName).Value.String() == "true" {
		a.list["metrics"] = a.createMetrics()
	}

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
	s := a.new(WithAddr(fmt.Sprintf(":%d", HealthzPort)))
	return s
}

// Function to create a new server for metrics
func (a *app) createMetrics() *server {
	s := a.new(WithAddr(fmt.Sprintf(":%d", metricsPort)))
	s.Config.Get(metricsPath, metrics.Handler().ServeHTTP)
	return s
}

// Function return a server
func (a *app) new(opts ...Option) *server {
	// create a new router
	r := chi.NewRouter()
	r.Use(logger.Logger("router", log.GetLogger()))
	r.Use(middleware.Recoverer)

	// create a new server with default parameters
	s := &server{
		http: &http.Server{
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
		panic(fmt.Sprintf("invalid address: %s", addr))
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

// Function WithCustomHandlerForHealth return a function Option
// Function take in parameter a function that return a boolean and an error
// and the endpoint path (e.g. /healthz)
func WithCustomHandlerForHealth(req HealthzFunc) OptionServer {
	return func(a *app) {
		// Prevent panic if the health server is not enabled
		if _, ok := a.list["health"]; !ok {
			return
		}

		a.list["health"].Config.Get(HealthzPath, func(w http.ResponseWriter, r *http.Request) {
			ok, err := req()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "application/json")
			if !ok {
				_, err = w.Write([]byte(`{"status":"ok"}`))
			} else {
				_, err = w.Write([]byte(`{"status":"ko"}`))
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	}
}
