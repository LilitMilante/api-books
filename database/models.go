package database

import "errors"

type Book struct {
	ID     int64  `json:"id,omitempty"`
	Title  string `json:"title,omitempty"`
	Author Author `json:"author,omitempty"`
}

type Author struct {
	ID        int64  `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
}

func (b *Book) Validate() error {
	if b.Title == "" {
		return errors.New("book title is required")
	}

	if err := b.Author.Validate(); err != nil {
		return err
	}

	return nil
}

func (a *Author) Validate() error {
	if a.Firstname == "" && a.Lastname == "" {
		return errors.New("firstname and lastname are required")
	}

	if a.Firstname == "" {
		return errors.New("author firstname is required")
	}

	if a.Lastname == "" {
		return errors.New("author lastname is required")
	}

	return nil
}
