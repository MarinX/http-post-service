package main

import (
	"io/ioutil"
	"net/http"
)

const (
	// StatusURLMissing error message for when URL is not specified
	StatusURLMissing = "Parametar URL is missing in body request"
)

// Server is custom wrapper for listen and serve
type Server struct {
	addr   string
	notify NotificationService
}

// NotificationService handles notification for server
type NotificationService interface {
	OnRequest(string) error
}

// NewServer creates a new custom server
func NewServer(addr string, notifyService NotificationService) *Server {
	if len(addr) == 0 {
		addr = ":8080"
	}
	return &Server{addr: addr, notify: notifyService}
}

// ListenAndServe is a wrapper for http.ListenAndServe
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.addr, s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	buff, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	if s.notify == nil {
		return
	}

	err = s.notify.OnRequest(string(buff))
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(err.Error()))
	}
}
