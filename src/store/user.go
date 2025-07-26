package store

import (
	"fmt"
	"time"

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

func (s *userStore) GetByName(name string) (domain.User, error) {
	var user domain.User

	err := s.db.Get(&user, "SELECT * FROM users WHERE name = $1", name)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (s *userStore) CreateUser(newUser domain.User) (domain.User, error) {
	if newUser.ID == uuid.Nil {
		newUser.ID = uuid.New()
	}

	now := time.Now()
	newUser.CreatedAt = now
	newUser.UpdatedAt = now

	_, err := s.db.Exec(
		`INSERT INTO users (id, name, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`,
		newUser.ID,
		newUser.Name,
		newUser.Password,
		newUser.CreatedAt,
		newUser.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

func (s *userStore) GetOne(id uuid.UUID) (domain.User, error) {
	var user domain.User

	err := s.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
