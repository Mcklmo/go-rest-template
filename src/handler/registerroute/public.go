package registerroute

import (
	"context"
	"fmt"
	"log/slog"

	"template/src/auth"
	"template/src/domain"

	"github.com/danielgtaylor/huma/v2"
)

func Public(api huma.API, userStore domain.UserStore, authHandler auth.AuthHandler, logger *slog.Logger) {
	register(
		api,
		makeOperation("POST", "/signup"),
		func(ctx context.Context, input *domain.CreateUserInput,
		) (*auth.SessionOutput, error) {
			user, err := domain.NewUser(input.Body.Name, input.Body.Password, userStore)
			if err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}

			session, err := authHandler.NewSession(user)
			if err != nil {
				return nil, err
			}

			return &session, nil
		},
		logger,
	)
	register(
		api,
		makeOperation("POST", "/login"),
		func(ctx context.Context, input *domain.CreateUserInput,
		) (*auth.SessionOutput, error) {
			user, err := domain.Login(input.Body.Name, input.Body.Password, userStore)
			if err != nil {
				return nil, fmt.Errorf("failed to login: %w", err)
			}

			session, err := authHandler.NewSession(user)
			if err != nil {
				return nil, err
			}

			return &session, nil
		},
		logger,
	)
}
