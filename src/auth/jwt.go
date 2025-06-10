package auth

import (
	"fmt"
	"time"

	"template/src/domain"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/lestrrat-go/jwx/v3/jwt/openid"
)

func NewAuthHandler(privateKey any, publicKey any) *AuthHandler {
	return &AuthHandler{privateKey: privateKey, PublicKey: publicKey}
}

type (
	SessionOutput struct {
		Body struct {
			Token string `json:"token"`
		}
	}
	AuthHandler struct {
		privateKey any
		PublicKey  any
	}
	Scope string
)

const (
	ScopeDomain Scope = "domain"
)

const (
	ScopesKey = "scopes"
	UserIDKey = "user_id"
)

func (h *AuthHandler) NewSession(user domain.User) (_ SessionOutput, err error) {
	var token openid.Token

	if token, err = makeToken(int64(3600), user); err != nil {
		return SessionOutput{}, fmt.Errorf("failed to generate token: %w", err)
	}

	tokenString, err := h.sign(token)
	if err != nil {
		return SessionOutput{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return SessionOutput{
		Body: struct {
			Token string `json:"token"`
		}{Token: tokenString},
	}, nil
}

func (h AuthHandler) sign(t openid.Token) (string, error) {
	signed, err := jwt.Sign(t, jwt.WithKey(jwa.RS256(), h.privateKey))
	if err != nil {
		return "", err
	}

	return string(signed), nil
}

func makeToken(
	expiresInSeconds int64,
	user domain.User,
) (t openid.Token, err error) {
	toSet := map[string]any{
		UserIDKey:       user.ID.String(),
		ScopesKey:       GetScopes(user),
		jwt.IssuedAtKey: time.Now(),
	}

	if expiresInSeconds != 0 {
		toSet[jwt.ExpirationKey] = time.Now().Add(time.Second * time.Duration(expiresInSeconds))
	}

	t = openid.New()

	for key, value := range toSet {
		if err = t.Set(key, value); err != nil {
			return nil, fmt.Errorf("failed to set %s: %w", key, err)
		}
	}

	return t, nil
}

func GetScopes(u domain.User) map[string]any {
	return map[string]any{
		string(ScopeDomain): true,
	}
}
