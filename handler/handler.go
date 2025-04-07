package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer(muxServer *mux.Router) *Server {
	s := &Server{
		router: muxServer,
	}
	s.router.Use(logRequest)

	s.router.HandleFunc("/", s.handlerHelloWorld).Methods("GET")

	return s
}

func (s *Server) handlerHelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("hello world"))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
