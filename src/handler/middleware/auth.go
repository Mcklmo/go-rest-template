package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"template/src/auth"
	"template/src/domain"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/lestrrat-go/jwx/v3/jwt/openid"
)

type (
	evaluatedToken struct {
		User domain.User
	}
)

func Auth(
	publicKey any,
	logger *slog.Logger,
	api huma.API,
	userStore domain.UserStore,
) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		allowingScopes := allowingScopes(ctx)
		if len(allowingScopes) == 0 {
			next(ctx)
			return
		}

		parsedJWT, err := parse(
			strings.TrimPrefix(ctx.Header("Authorization"), "Bearer "),
			publicKey,
			allowingScopes,
		)
		if err != nil {
			unauthorized(ctx, fmt.Sprintf("error parsing jwt: %v", err), api, logger)
			return
		}

		token, err := EvaluateToken(parsedJWT, userStore)
		if err != nil {
			unauthorized(ctx, fmt.Sprintf("error evaluating token: %v", err), api, logger)
			return
		}

		if !anyScopeMatches(allowingScopes, auth.GetScopes(token.User)) {
			unauthorized(ctx, "user does not have any allowed scopes", api, logger)
			return
		}

		ctx = addValue(ctx, token.User)

		next(ctx)
	}
}

func EvaluateToken(parsedToken jwt.Token, userStore domain.UserStore) (r evaluatedToken, err error) {
	user, err := getUserFromToken(parsedToken, userStore)
	if err != nil {
		return evaluatedToken{}, err
	}

	return evaluatedToken{User: user}, nil
}

func getUserFromToken(parsedToken jwt.Token, userStore domain.UserStore) (user domain.User, err error) {
	var userID_ string

	if err = parsedToken.Get(auth.UserIDKey, &userID_); err != nil {
		return domain.User{}, fmt.Errorf("no user found in token: %w", err)
	}

	userID, err := uuid.Parse(userID_)
	if err != nil {
		return domain.User{}, fmt.Errorf("error parsing user id: %w", err)
	}

	user, err = userStore.GetOne(userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("error getting user[%s] from db: %w", userID_, err)
	}

	return user, nil
}

func addValue(ctx huma.Context, val domain.Contextable) huma.Context {
	return huma.WithValue(ctx, val.GetContextKey(), val.GetValue())
}

func parse(rawToken string, publicKey any, allowingScopes []string) (jwt.Token, error) {
	if len(rawToken) == 0 {
		return nil, errors.New("no token found")
	}

	parsedToken, err := ParseToken(rawToken, publicKey)
	if err != nil {
		return nil, fmt.Errorf("token invalid: %w", err)
	}

	if parsedToken == nil {
		return nil, fmt.Errorf("token parsing failed")
	}

	if err = hasAnyAllowedScope(parsedToken, allowingScopes); err != nil {
		return nil, err
	}

	return parsedToken, nil
}

func hasAnyAllowedScope(parsedToken jwt.Token, allowedScopes []string) (err error) {
	var userScopes map[string]any

	if err = parsedToken.Get(auth.ScopesKey, &userScopes); err != nil {
		return fmt.Errorf("error getting scopes from token: %w", err)
	}

	if anyScopeMatches(allowedScopes, userScopes) {
		return nil
	}

	userScopesString := ""

	for scope, ok := range userScopes {
		_ok := ok.(bool)
		if !_ok {
			continue
		}

		userScopesString += scope + ", "
	}

	userScopesString = strings.TrimSuffix(userScopesString, ", ")

	return fmt.Errorf("user has scopes %#v, expected one of: %v", userScopesString, allowedScopes)
}

func anyScopeMatches(allowedScopes []string, userScopes map[string]any) bool {
	for _, allowedScope := range allowedScopes {
		userScope, ok := userScopes[allowedScope]
		if !ok {
			continue
		}

		if valid, ok := userScope.(bool); ok && valid {
			return true
		}
	}

	return false
}

// allowingScopes gets the scopes that are allowed to access the operation. If none are found, the operation is not protected.
func allowingScopes(ctx huma.Context) []string {
	var scopesThatAllowAccess []string

	for _, opScheme := range ctx.Operation().Security {
		var ok bool

		if scopesThatAllowAccess, ok = opScheme["scope"]; ok {
			break
		}
	}

	return scopesThatAllowAccess
}

func unauthorized(ctx huma.Context, msg string, api huma.API, logger *slog.Logger) {
	logger.Error(msg)

	err := huma.WriteErr(api, ctx, http.StatusUnauthorized, strings.Split(msg, ":")[0])
	if err != nil {
		logger.Error("failed to write huma error", "error", err)
	}
}

func ParseToken(rawToken string, publicKey any) (jwt.Token, error) {
	return jwt.Parse([]byte(rawToken), getJWTParsingOptions(publicKey)...)
}

func getJWTParsingOptions(publicKey any) []jwt.ParseOption {
	return []jwt.ParseOption{jwt.WithKey(jwa.RS256(), publicKey), jwt.WithToken(openid.New())}
}

func ParseTokenInsecure(rawToken string, publicKey any) (jwt.Token, error) {
	return jwt.ParseInsecure([]byte(rawToken), getJWTParsingOptions(publicKey)...)
}
