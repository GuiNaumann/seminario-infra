
package repository

import (
	"database/sql"

	"biblioteca-api/models"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) GetAll() ([]models.Book, error) {
	rows, err := r.db.Query(`SELECT id, title, author, year, available FROM books ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.Available); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (r *BookRepository) GetByID(id int) (*models.Book, error) {
	var b models.Book
	err := r.db.QueryRow(
		`SELECT id, title, author, year, available FROM books WHERE id = $1`, id,
	).Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.Available)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &b, err
}

func (r *BookRepository) Create(b *models.Book) error {
	return r.db.QueryRow(
		`INSERT INTO books (title, author, year, available) VALUES ($1, $2, $3, $4) RETURNING id`,
		b.Title, b.Author, b.Year, b.Available,
	).Scan(&b.ID)
}

func (r *BookRepository) Update(id int, b *models.Book) (bool, error) {
	res, err := r.db.Exec(
		`UPDATE books SET title=$1, author=$2, year=$3, available=$4 WHERE id=$5`,
		b.Title, b.Author, b.Year, b.Available, id,
	)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (r *BookRepository) Delete(id int) (bool, error) {
	res, err := r.db.Exec(`DELETE FROM books WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}
package repository

import (
	"database/sql"

	"biblioteca-api/models"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) GetAll() ([]models.Book, error) {
	rows, err := r.db.Query(`SELECT id, title, author, year, available FROM books ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.Available); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (r *BookRepository) GetByID(id int) (*models.Book, error) {
	var b models.Book
	err := r.db.QueryRow(
		`SELECT id, title, author, year, available FROM books WHERE id = $1`, id,
	).Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.Available)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &b, err
}

func (r *BookRepository) Create(b *models.Book) error {
	return r.db.QueryRow(
		`INSERT INTO books (title, author, year, available) VALUES ($1, $2, $3, $4) RETURNING id`,
		b.Title, b.Author, b.Year, b.Available,
	).Scan(&b.ID)
}

func (r *BookRepository) Update(id int, b *models.Book) (bool, error) {
	res, err := r.db.Exec(
		`UPDATE books SET title=$1, author=$2, year=$3, available=$4 WHERE id=$5`,
		b.Title, b.Author, b.Year, b.Available, id,
	)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (r *BookRepository) Delete(id int) (bool, error) {
	res, err := r.db.Exec(`DELETE FROM books WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}
