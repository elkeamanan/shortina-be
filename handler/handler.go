package handler

import (
	linkService "elkeamanan/shortina/internal/link/service"
	userService "elkeamanan/shortina/internal/user/service"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router      *mux.Router
	linkService linkService.LinkService
	userService userService.UserService
}

func NewServer(muxServer *mux.Router, linkService linkService.LinkService, userService userService.UserService) *Server {
	s := &Server{
		router:      muxServer,
		linkService: linkService,
		userService: userService,
	}

	s.router.Use(logRequest)
	s.router.Use(enableCORS)

	// public APIs
	s.router.HandleFunc("/", s.handlerHelloWorld).Methods("GET")
	s.router.HandleFunc("/links/{key}", s.handlerGetLinkRedirection).Methods("GET")
	s.router.HandleFunc("/users/register", s.handlerRegisterUser).Methods("POST")
	s.router.HandleFunc("/users/login", s.handlerLoginUser).Methods("POST")
	s.router.HandleFunc("/refresh-token", s.handlerRefreshToken).Methods("GET")

	// optional auth
	optionalAuthRoutes := s.router.NewRoute().Subrouter()
	optionalAuthRoutes.Use(optionalAuthVerifier)
	optionalAuthRoutes.HandleFunc("/links", s.handlerStoreLink).Methods("POST", "OPTIONS")

	// requires auth
	protectedRoutes := s.router.NewRoute().Subrouter()
	protectedRoutes.Use(authVerifier)
	protectedRoutes.HandleFunc("/users/{id}/update", s.handlerUpdateUser).Methods("PATCH")

	return s
}

func (s *Server) handlerHelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("hello world"))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
