package api

import (
	"apiBooks/database"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) getBooks(w http.ResponseWriter, _ *http.Request) {
	books, err := s.store.Books()
	if err != nil {
		log.Println("Get books:", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		log.Println("Encode books:", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (s *Server) createBook(w http.ResponseWriter, r *http.Request) {
	var book database.Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Println("Decode book:", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if err = book.Validate(); err != nil {
		log.Println("Book validation:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	author, err := s.store.AuthorByFullName(book.Author.Firstname, book.Author.Lastname)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Get author:", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if author.ID == 0 {
		book.Author.ID, err = s.store.AddAuthor(book.Author)
		if err != nil {
			log.Println("Create author:", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	} else {
		book.Author.ID = author.ID
	}

	id, err := s.store.AddBook(book)
	if err != nil {
		log.Println("Add book to db:", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	book.ID = id

	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		log.Println("Encode book:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) createBooks(w http.ResponseWriter, r *http.Request) {
	books := make([]database.Book, 0)

	if err := json.NewDecoder(r.Body).Decode(&books); err != nil {
		log.Println("Decode books:", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := s.store.AddBooks(books); err != nil {
		log.Println("Add books:", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if err := json.NewEncoder(w).Encode(books); err != nil {
		log.Println("Encode books:", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) updateBook(w http.ResponseWriter, r *http.Request) {
	var book database.Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Println("Decode book:", err)
		w.WriteHeader(http.StatusBadRequest)

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

	book.ID = int64(id)

	err = s.store.UpdateBooks(book)
	if err != nil {
		log.Println("Update book:", err)
		w.WriteHeader(http.StatusNotFound)

		return
	}

	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		log.Println("Encode book:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Convert to int", idStr, err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	bookID := int64(id)

	err = s.store.DeleteBook(bookID)
	if err != nil {
		log.Println("Delete book:", err)
		w.WriteHeader(http.StatusNotFound)

		return
	}
}
