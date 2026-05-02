package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	config *Config
	router chi.Router
}

func NewServer(config *Config) *Server {
	return &Server{
		config: config,
		router: chi.NewRouter(),
	}
}

func (s *Server) Start() error {
	s.configureRouter()

	return http.ListenAndServe(s.config.Addr, s.router)
}

func (s *Server) configureRouter() {
	s.router.Get("/", s.HandleRoot())
}

func (s *Server) HandleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		http.ServeFile(w, request, "./web/index.html")
	}
}
