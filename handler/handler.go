package handler

import (
	linkService "elkeamanan/shortina/internal/link/service"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router      *mux.Router
	linkService linkService.LinkService
}

func NewServer(muxServer *mux.Router, linkService linkService.LinkService) *Server {
	s := &Server{
		router:      muxServer,
		linkService: linkService,
	}

	s.router.Use(logRequest)
	s.router.Use(enableCORS)

	s.router.HandleFunc("/", s.handlerHelloWorld).Methods("GET")
	s.router.HandleFunc("/link", s.handlerStoreLink).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/link/{key}", s.handlerGetLinkRedirection).Methods("GET")

	return s
}

func (s *Server) handlerHelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("hello world"))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
