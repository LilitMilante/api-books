package api

import (
	"apiBooks/database"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) getAuthors(w http.ResponseWriter, _ *http.Request) {
	authors, err := s.store.Authors()
	if err != nil {
		log.Println("Get authors:", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(authors); err != nil {
		log.Println("Encode authors:", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (s *Server) createAuthor(w http.ResponseWriter, r *http.Request) {
	var author database.Author

	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		log.Println("Decode author:", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := author.Validate(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	id, err := s.store.AddAuthor(author)
	if err != nil {
		log.Println("Add author to db:", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	author.ID = id

	if err := json.NewEncoder(w).Encode(author); err != nil {
		log.Println("Encode author:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (s *Server) updateAuthor(w http.ResponseWriter, r *http.Request) {
	var author database.Author

	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		log.Println("Decode author:", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := author.Validate(); err != nil {
		log.Println("Author update:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Convert to int", idStr, err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	author.ID = int64(id)

	if err = s.store.UpdateAuthor(author); err != nil {
		log.Println("Update book:", err)
		w.WriteHeader(http.StatusNotFound)

		return
	}

	if err = json.NewEncoder(w).Encode(author); err != nil {
		log.Println("Encode book:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) deleteAuthor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Convert to int", idStr, err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	authorID := int64(id)

	if err = s.store.DeleteAuthor(authorID); err != nil {
		log.Println("Delete author:", err)
		w.WriteHeader(http.StatusNotFound)

		return
	}

}
