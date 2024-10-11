package health

// import (
// 	"time"
// )

// const (
// 	timeoutR = 1 * time.Second
// )

// healthHandler returns a http.Handler that returns a health check response
// func DefaultHandler() http.Handler {
// 	// TODO - Implement a new way to ask the health of the application (e.g. check image updater)
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		_, err := net.DialTimeout("tcp", ":9081", timeoutR)
// 		if err != nil {
// 			return
// 		}

// 		// TODO - Implement an http.Handler content-type
// 		w.Header().Set("Content-Type", "application/json")
// 		_, err = w.Write([]byte(`{"status":"ok"}`))
// 		if err != nil {
// 			return
// 		}
// 	})
// }
