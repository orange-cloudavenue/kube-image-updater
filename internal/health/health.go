package health

import (
	"context"
	"flag"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/orange-cloudavenue/kube-image-updater/internal/httpserver"
)

const (
	timeoutR = 1 * time.Second
)

var (
	healthPath string = "/healthz"
	healthPort string = ":9081"
)

func init() {
	flag.StringVar(&healthPort, "health-port", healthPort, "Health server port. ex: :9081")
	flag.StringVar(&healthPath, "health-path", healthPath, "Health server path. ex: /healthz")
}

// healthHandler returns a http.Handler that returns a health check response
func healthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := net.DialTimeout("tcp", healthPort, timeoutR)
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(`{"status":"ok"}`))
		if err != nil {
			return
		}
	})
}

// ServeHealth starts the health check server
func StartHealth(ctx context.Context, wg *sync.WaitGroup) (err error) {
	s := httpserver.New(httpserver.WithAddr(healthPort))
	s.AddGetRoutes(healthPath, healthHandler())
	return s.Start(ctx, wg)
}
