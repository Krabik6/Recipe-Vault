package repository

import "database/sql"

type Authorization struct {
}

type Repository struct {
	Authorization
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: Authorization{},
	}
}
