package health

import (
	"context"
	"errors"
	"flag"
	"log"
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

// ServeHealth starts the health check server
func ServeHealth(ctx context.Context) (err error) {
	// Define Health check server
	mux := http.NewServeMux()
	mux.HandleFunc(healthPath, func(w http.ResponseWriter, r *http.Request) {
		// TODO - Add more health checks like use of kube client on kube api server
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

	// create health check server
	s := &http.Server{
		Addr:        healthPort,
		Handler:     mux,
		ReadTimeout: 10 * timeoutR,
	}

	// start the HTTP server
	go func() {
		log.Printf("Starting health check server on %s", s.Addr)
		if err = s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return
		}
	}()

	// kill the server if there is an error
	go func() {
		for {
			<-ctx.Done()
			ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
			log.Printf("Shutting down health check server on %s", s.Addr)
			defer cancel()
			if err = s.Shutdown(ctxTimeout); err != nil {
				log.Printf("Failed to shutdown health check server: %v", err)
			}
			return
		}
	}()

	return nil
}
