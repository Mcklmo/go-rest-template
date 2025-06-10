package registerroute

import (
	"context"
	"log/slog"
	"net/http"

	"template/src/auth"
	"template/src/domain"

	"github.com/danielgtaylor/huma/v2"
)

func Domain(api huma.API, userStore domain.UserStore, store domain.Store, authHandler *auth.AuthHandler, logger *slog.Logger) {
	register(
		api,
		makeSecuredOperation(http.MethodGet, "/route", auth.ScopeDomain),
		func(ctx context.Context, input *struct{},
		) (*domain.Output, error) {
			user, err := domain.GetCTX[domain.User](ctx)
			if err != nil {
				return nil, err
			}

			_ = user

			return &domain.Output{}, nil
		},
		logger,
	)
	register(
		api,
		makeSecuredOperation(http.MethodPost, "/route", auth.ScopeDomain),
		func(ctx context.Context, input *domain.Input,
		) (*domain.Output, error) {
			user, err := domain.GetCTX[domain.User](ctx)
			if err != nil {
				return nil, err
			}

			_ = user

			return &domain.Output{}, nil
		},
		logger,
	)
}
