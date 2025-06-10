package registerroute

import (
	"context"
	"log/slog"

	"template/src/auth"
	"template/src/domain"

	"github.com/danielgtaylor/huma/v2"
)

func Public(api huma.API, userStore domain.UserStore, authHandler *auth.AuthHandler, logger *slog.Logger) {
	register(
		api,
		makeOperation("POST", "/signup"),
		func(ctx context.Context, input *domain.CreateUserInput,
		) (*auth.SessionOutput, error) {
			user, err := domain.NewUser(input.Body.Name, userStore)
			if err != nil {
				return nil, err
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
