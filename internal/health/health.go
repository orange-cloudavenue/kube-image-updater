package health

import (
	"flag"
	"net"
	"net/http"
	"time"
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
func Handler() http.Handler {
	// TODO - Implement a new way to ask the health of the application (e.g. check image updater)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := net.DialTimeout("tcp", healthPort, timeoutR)
		if err != nil {
			return
		}

		// TODO - Implement an http.Handler content-type
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(`{"status":"ok"}`))
		if err != nil {
			return
		}
	})
}
