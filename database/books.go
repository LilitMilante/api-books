package database

import (
	"database/sql"
	"errors"
)

func (s *Storage) Books() ([]Book, error) {
	books := make([]Book, 0)

	rows, err := s.conn.Query(
		`
			select 
			b.id b_id, 
			b.title b_title, 
			a.id a_id, 
			a.firstname, 
			a.lastname
			from books b
			join authors a
			on a.id = b.author_id 
			`,
	)
	if err != nil {
		return books, err
	}

	defer rows.Close()

	for rows.Next() {
		var book Book

		if err := rows.Scan(&book.ID, &book.Title, &book.Author.ID, &book.Author.Firstname, &book.Author.Lastname); err != nil {
			break
		}

		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return books, err
	}

	return books, nil
}

func (s *Storage) AddBook(b Book) (int64, error) {
	var id int64

	result := s.conn.QueryRow("insert into books (title, author_id) values($1, $2) returning id", b.Title, b.Author.ID)

	if err := result.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) UpdateBooks(b Book) error {
	result, err := s.conn.Exec("update books set title = $1, author_id = $2 where id = $3", b.Title, b.Author, b.ID)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count != 1 {
		return errors.New("book is not found")
	}

	return nil
}

func (s *Storage) DeleteBook(id int64) error {
	result, err := s.conn.Exec("delete from books where id = $1", id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count != 1 {
		return errors.New("book is not found")
	}

	return nil
}

func (s *Storage) AddBooks(books []Book) error {
	stmt, err := s.conn.Prepare("insert into books (title, author_id) values($1, $2) returning id")
	if err != nil {
		return err
	}

	defer stmt.Close()

	for i, book := range books {
		row := stmt.QueryRow(book.Title, book.Author)

		var id int64
		if err := row.Scan(&id); err != nil {
			return err
		}

		books[i].ID = id
	}

	return nil
}

func (s *Storage) BooksByAuthorId(id int64) ([]Book, error) {
	books := make([]Book, 0)

	rows, err := s.conn.Query(
		`
				select b.id, 
				       b.title, 
				       a.id, 
				       a.firstname, 
				       a.lastname
				from books b 
				join authors a on a.id = b.author_id
				where a.id = $1
				       `,
		id,
	)

	if err != nil {
		if err != sql.ErrNoRows {
			return []Book{}, err
		}

		return books, nil
	}

	defer rows.Close()

	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author.ID, &b.Author.Firstname, &b.Author.Lastname); err != nil {
			break
		}

		books = append(books, b)
	}

	if err = rows.Err(); err != nil {
		return []Book{}, err
	}

	return books, nil
}
