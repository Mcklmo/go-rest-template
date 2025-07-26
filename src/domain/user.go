package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	User struct {
		ID        uuid.UUID `db:"id"`
		Name      string    `db:"name"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
		Password  string    `db:"password"`
	}
	CreateUserInput struct {
		Body struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}
	}
	UserStore interface {
		CreateUser(newUser User) (User, error)
		GetOne(id uuid.UUID) (User, error)
		GetByName(name string) (User, error)
	}
	// ContextKey is a custom type to avoid key collisions in context values.
	ContextKey  string
	Contextable interface {
		GetContextKey() ContextKey
		GetValue() any
	}
)

const UserContextKey ContextKey = "user"

func GetCTX[T Contextable](ctx context.Context) (r T, err error) {
	var ok bool

	k := r.GetContextKey()
	v := ctx.Value(k)

	if r, ok = v.(T); !ok {
		err = huma.Error500InternalServerError(fmt.Sprintf("checking context for '%s', wanted %T, got %T", k, r, v))
		return
	}

	return
}

func NewUser(name string, password string, store UserStore) (user User, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	if user, err = store.CreateUser(User{Name: name, Password: string(hashedPassword)}); err != nil {
		return User{}, err
	}

	return user, nil
}

func (u User) GetContextKey() ContextKey {
	return UserContextKey
}

func (u User) GetValue() any {
	return u
}

func Login(name string, password string, userStore UserStore) (User, error) {
	user, err := userStore.GetByName(name)
	if err != nil {
		return User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return User{}, err
	}

	return user, nil
}
