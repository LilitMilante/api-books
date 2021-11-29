package api

import (
	"apiBooks/database"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	port   string
	router *mux.Router
	store  *database.Storage
}

func NewServer(p string, r *mux.Router, s *database.Storage) *Server {
	return &Server{
		port:   p,
		router: r,
		store:  s,
	}
}

func (s *Server) cors(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")

		if r.Method == http.MethodOptions {
			return
		}

		nextHandler.ServeHTTP(w, r)
	})
}

func (s *Server) Start() error {
	s.router.Use(s.cors)

	s.router.HandleFunc("/books", s.getBooks).Methods(http.MethodGet).Queries("aid", "{aid:[0-9]+}")
	s.router.HandleFunc("/books", s.createBook).Methods(http.MethodPost)
	s.router.HandleFunc("/books/{id}", s.updateBook).Methods(http.MethodPut)
	s.router.HandleFunc("/books/{id}", s.deleteBook).Methods(http.MethodDelete)
	s.router.HandleFunc("/books/many", s.createBooks).Methods(http.MethodPost)

	s.router.HandleFunc("/authors", s.getAuthors).Methods(http.MethodGet)
	s.router.HandleFunc("/authors", s.createAuthor).Methods(http.MethodPost)
	s.router.HandleFunc("/authors/{id}", s.updateAuthor).Methods(http.MethodPut)
	s.router.HandleFunc("/authors/{id}", s.deleteAuthor).Methods(http.MethodDelete)

	return http.ListenAndServe(s.port, s.router)
}
