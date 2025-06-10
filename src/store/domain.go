package store

import (
	"template/src/domain"

	"github.com/jmoiron/sqlx"
)

type store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) domain.Store {
	return &store{db: db}
}
