package service

import (
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
}

type Service struct {
	Authorization
}

//repos *Repository.Repository

func NewService(db *sqlx.DB) *Service {
	return &Service{}
}
