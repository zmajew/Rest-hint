package server

import (
	"Rest-hint/config"
	"Rest-hint/service"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	config       *config.Config
	movieService *service.MovieService
}

func (s *Server) Handler() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hint/{lett}", s.movies).Methods("GET")
	return r
}

func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.Handler(),
	}

	httpServer.ListenAndServe()
}

func (s *Server) movies(w http.ResponseWriter, r *http.Request) {
	letters := mux.Vars(r)

	movies, err := s.movieService.GiveHint(letters)
	if err != nil {
		log.Println("error 3 quering database:", err)
		w.WriteHeader(500)
		w.Write([]byte("error 4 quering database:"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(movies)
	if err != nil {
		log.Println("error writing search response:", err)
		w.WriteHeader(500)
		w.Write([]byte("An error occurred writing response"))
	}
}

func NewServer(config *config.Config, service *service.MovieService) *Server {
	return &Server{
		config:       config,
		movieService: service,
	}
}
