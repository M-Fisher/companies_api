package server

import (
	"encoding/json"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func (s *Server) initApp() *http.Server {
	r := s.createHTTPHandler()

	http.Handle("/", r)

	return &http.Server{
		Addr:    s.Config.Port,
		Handler: r,
	}
}

func (s *Server) CreateDefaultRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", s.createHealthChecker())
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/{action}", pprof.Index)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.Router = r
	return s.Router
}

// CreateHTTPHandler creates handler for API
func (s *Server) createHTTPHandler() http.Handler {
	headersOk := handlers.AllowedHeaders([]string{
		"X-Requested-With",
		"Authorization",
		"authorization",
		"sentry-trace",
		"content-type",
		"Content-Type",
	})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	originsOk := handlers.AllowedOrigins([]string{"*"})

	return handlers.CORS(originsOk, headersOk, methodsOk)(s.Router)
}

func (s *Server) createHealthChecker() func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		type ServiceStatus struct {
			Status string `json:"status"`
		}
		type HealthCheckResponse struct {
			Timestamp string                   `json:"timestamp"`
			Services  map[string]ServiceStatus `json:"service,omitempty"`
			Errors    map[string]string        `json:"errors,omitempty"`
		}
		ErrorsList := map[string]string{}
		ServicesList := map[string]ServiceStatus{}
		err := s.Storage.GetStatus()
		if err != nil {
			ErrorsList["DB"] = err.Error()
		} else {
			ServicesList["DB"] = ServiceStatus{
				Status: `OK`,
			}
		}

		res, _ := json.Marshal(HealthCheckResponse{
			Timestamp: time.Now().Format(time.RFC3339),
			Services:  ServicesList,
			Errors:    ErrorsList,
		})

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write(res)
	}
}
