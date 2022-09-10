package base

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/M-Fisher/companies_api/app/internal/logger"
	"github.com/M-Fisher/companies_api/app/internal/server"
	"github.com/M-Fisher/companies_api/app/internal/services/auth"
)

type APIInterface interface {
	SetRouter(*mux.Router)
	SetRoutes()
}

type Response map[string]interface{}

// API all other API should extend this struct
type API struct {
	Srv    *server.Server
	Router *mux.Router
}

func (a *API) VerifyUser(r *http.Request) (uint64, error) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	if len(splitToken) > 1 {
		user, err := new(auth.JWTUser).Parse(strings.TrimSpace(splitToken[1]), a.Srv.Config.JWTSecret)
		if err == nil {
			return uint64(user.ID), nil
		}
		return 0, ErrUnauthorized
	}
	return 0, ErrUnauthorized
}

func (a *API) SetJSONHandler(path string, f func(rw http.ResponseWriter, r *http.Request) (interface{}, error)) *mux.Route {
	if a.Router != nil {
		return a.Router.HandleFunc(path,
			a.loggingMiddleware(
				a.panicHandlerMiddleware(
					func(rw http.ResponseWriter, rq *http.Request) {
						var (
							res  any
							err  error
							done = make(chan struct{})
						)
						go func() {
							res, err = f(rw, rq)
							done <- struct{}{}
						}()
						for {
							select {
							case <-rq.Context().Done():
								// Received Done signal from parent
								a.sendErrorResponse(rw, errors.New("request cancelled"))
								return
							case <-done:
								if err != nil {
									errJS, _ := json.Marshal(err)
									a.Srv.Log.Debug("Outcoming Response",
										zap.String("status", "err"),
										zap.String("error", string(errJS)),
									)
									a.sendErrorResponse(rw, err)
								} else {
									resJS, _ := json.Marshal(res)
									a.Srv.Log.Debug("Outcoming Response",
										zap.String("status", "ok"),
										zap.String("res", string(resJS)),
									)
									a.sendJSONResponse(rw, res)
								}
								return
							}
						}

					},
				),
			),
		)
	} else {
		return nil
	}
}

func (a *API) panicHandlerMiddleware(
	next func(rw http.ResponseWriter, rq *http.Request),
) func(rw http.ResponseWriter, rq *http.Request) {
	return func(rw http.ResponseWriter, rq *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				s := debug.Stack()
				a.Srv.Log.Error("Got panic in api JsonHandler", zap.String("stack", string(s)))
				a.sendErrorResponse(rw, errors.New(`internal server error`))
			}
		}()
		next(rw, rq)
	}
}

func (a *API) loggingMiddleware(
	next func(rw http.ResponseWriter, rq *http.Request),
) func(rw http.ResponseWriter, rq *http.Request) {
	return func(rw http.ResponseWriter, rq *http.Request) {
		log := a.Srv.Log.
			With(zap.String("request_uri", rq.RequestURI),
				zap.String("uri", rq.RequestURI),
				zap.String("host", rq.Host),
				zap.String("http_method", rq.Method),
				zap.String("trace_id", CreateGUID()),
			)
		headers := &bytes.Buffer{}
		for k, v := range rq.Header {
			headers.WriteString(fmt.Sprintf("%s:%s,", k, v))
		}
		remoteAddr := rq.RemoteAddr
		if chunks := strings.Split(remoteAddr, ":"); len(chunks) > 0 {
			remoteAddr = chunks[0]
		}
		log.Debug("Incoming Request",
			zap.String("headers", headers.String()),
			zap.String("remote_addr", remoteAddr),
		)

		rq = rq.WithContext(logger.ToContext(rq.Context(), log))

		next(rw, rq)
	}
}

func (a *API) sendJSONResponse(w http.ResponseWriter, payload any) {
	response := map[string]any{
		"status_code": 0,
	}
	if payload != nil {
		response["payload"] = payload
	}
	a.sendResponse(w, response, 0)
}

func (a *API) sendErrorResponse(w http.ResponseWriter, err error) {
	response := map[string]any{
		"status_code": ServerErrorCode,
		"status_text": err.Error(),
	}
	errorCode := http.StatusInternalServerError
	if errors.Is(err, ErrUnauthorized) {
		errorCode = http.StatusUnauthorized
	}
	a.sendResponse(w, response, errorCode)
}

func (a *API) sendResponse(w http.ResponseWriter, response map[string]any, errorCode int) {
	js, err := json.Marshal(response)

	if err != nil {
		a.Srv.Log.Error(`response marshalling error`, zap.Error(err))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Encoding", "gzip")

	if errorCode != 0 {
		w.WriteHeader(errorCode)
	}

	writer, err := gzip.NewWriterLevel(w, gzip.BestCompression)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	defer writer.Close()
	_, err = writer.Write(js)
	if err != nil {
		a.Srv.Log.Error("Failed to write response", zap.Error(err))
	}
}

func (a *API) SetRouter(r *mux.Router) {
	a.Router = r
}
