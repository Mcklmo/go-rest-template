package store

import (
	"template/src/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) domain.UserStore {
	return &userStore{db: db}
}

func (s *userStore) CreateUser(newUser domain.User) (domain.User, error) {
	if newUser.ID == uuid.Nil {
		newUser.ID = uuid.New()
	}

	err := s.db.Get(
		&newUser.ID,
		`INSERT INTO users (id, name)
		VALUES ($1, $2)
		RETURNING id`,
		newUser.ID,
		newUser.Name,
	)
	if err != nil {
		return domain.User{}, err
	}

	return newUser, nil
}

func (s *userStore) GetOne(id uuid.UUID) (domain.User, error) {
	var user domain.User

	err := s.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
