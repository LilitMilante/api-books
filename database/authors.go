package database

import "errors"

func (s *Storage) AuthorByFullName(firstName, lastName string) (Author, error) {
	var author Author

	row := s.conn.QueryRow("select id, firstname, lastname from authors where firstname = $1 and lastname = $2", firstName, lastName)

	if err := row.Scan(&author.ID, &author.Firstname, &author.Lastname); err != nil {
		return Author{}, err
	}

	return author, nil
}

func (s *Storage) Authors() ([]Author, error) {
	var authors []Author

	rows, err := s.conn.Query("select id, firstname, lastname from authors")
	if err != nil {
		return authors, err
	}

	defer rows.Close()

	for rows.Next() {
		var author Author

		if err := rows.Scan(&author.ID, &author.Firstname, &author.Lastname); err != nil {
			break
		}

		authors = append(authors, author)
	}

	if err = rows.Err(); err != nil {
		return authors, err
	}

	return authors, nil
}

func (s *Storage) AddAuthor(a Author) (int64, error) {
	var id int64

	result := s.conn.QueryRow("insert into authors (firstname , lastname) values($1, $2) returning id", a.Firstname, a.Lastname)

	if err := result.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) UpdateAuthor(a Author) error {
	result, err := s.conn.Exec("update authors set firstname = $1, lastname = $2 where id = $3", a.Firstname, a.Lastname, a.ID)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count != 1 {
		return errors.New("author is not found")
	}

	return nil
}

func (s *Storage) DeleteAuthor(id int64) error {
	result, err := s.conn.Exec("delete from authors where id = $1", id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count != 1 {
		return errors.New("author is not found")
	}

	return nil
}
