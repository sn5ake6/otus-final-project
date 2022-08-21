package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sn5ake6/otus-final-project/internal/storage"
)

const (
	authorizeURI = "/authorize/"
	resetURI     = "/reset/"
	blacklistURI = "/blacklist/"
	whitelistURI = "/whitelist/"
)

type Server struct {
	addr       string
	logger     Logger
	httpServer *http.Server
}

type Logger interface {
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
	LogHTTPRequest(r *http.Request, statusCode int, requestDuration time.Duration)
}

type Application interface {
	Authorize(ctx context.Context, authorize storage.Authorize) (bool, error)
	Reset(ctx context.Context, authorize storage.Authorize)
	AddToBlacklist(ctx context.Context, subnet string) error
	DeleteFromBlacklist(ctx context.Context, subnet string) error
	FindIPInBlacklist(ctx context.Context, ip string) (bool, error)
	AddToWhitelist(ctx context.Context, subnet string) error
	DeleteFromWhitelist(ctx context.Context, subnet string) error
	FindIPInWhitelist(ctx context.Context, ip string) (bool, error)
}

func NewRouter(logger Logger, app Application) http.Handler {
	handler := &HTTPHandler{app: app}
	router := mux.NewRouter()
	router.StrictSlash(true)

	router.HandleFunc(authorizeURI, loggingMiddleware(handler.Authorize, logger)).Methods(http.MethodPost)
	router.HandleFunc(resetURI, loggingMiddleware(handler.Reset, logger)).Methods(http.MethodPost)
	router.HandleFunc(blacklistURI, loggingMiddleware(handler.AddToBlacklist, logger)).Methods(http.MethodPost)
	router.HandleFunc(blacklistURI, loggingMiddleware(handler.DeleteFromBlacklist, logger)).Methods(http.MethodDelete)
	router.HandleFunc(whitelistURI, loggingMiddleware(handler.AddToWhitelist, logger)).Methods(http.MethodPost)
	router.HandleFunc(whitelistURI, loggingMiddleware(handler.DeleteFromWhitelist, logger)).Methods(http.MethodDelete)

	return router
}

func NewServer(addr string, logger Logger, app Application) *Server {
	s := &Server{
		addr:   addr,
		logger: logger,
	}

	httpServer := &http.Server{
		Addr:    addr,
		Handler: NewRouter(logger, app),
	}

	s.httpServer = httpServer

	return s
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info(fmt.Sprintf("HTTP server started: %s", s.addr))
	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info(fmt.Sprintf("HTTP server stopped: %s", s.addr))

	return s.httpServer.Shutdown(ctx)
}

type HTTPHandler struct {
	app Application
}

func (h *HTTPHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	authorize, err := h.getAutorize(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	res, err := h.app.Authorize(r.Context(), authorize)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")

	if res {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
		return
	}

	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte(`{"ok": false}`))
}

func (h *HTTPHandler) Reset(w http.ResponseWriter, r *http.Request) {
	authorize, err := h.getAutorize(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	h.app.Reset(r.Context(), authorize)
	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) AddToBlacklist(w http.ResponseWriter, r *http.Request) {
	subnet, err := h.getSubnet(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = h.app.AddToBlacklist(r.Context(), subnet)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *HTTPHandler) DeleteFromBlacklist(w http.ResponseWriter, r *http.Request) {
	subnet, err := h.getSubnet(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = h.app.DeleteFromBlacklist(r.Context(), subnet)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) AddToWhitelist(w http.ResponseWriter, r *http.Request) {
	subnet, err := h.getSubnet(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = h.app.AddToWhitelist(r.Context(), subnet)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *HTTPHandler) DeleteFromWhitelist(w http.ResponseWriter, r *http.Request) {
	subnet, err := h.getSubnet(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = h.app.DeleteFromWhitelist(r.Context(), subnet)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) getAutorize(r *http.Request) (storage.Authorize, error) {
	var authorize storage.Authorize

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return authorize, err
	}

	err = json.Unmarshal(data, &authorize)
	if err != nil {
		return authorize, err
	}

	return authorize, nil
}

func (h *HTTPHandler) getSubnet(r *http.Request) (string, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	var subnet storage.Subnet
	err = json.Unmarshal(data, &subnet)
	if err != nil {
		return "", err
	}

	return subnet.Subnet, nil
}
