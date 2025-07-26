package registerroute

import (
	"context"
	"log/slog"
	"net/http"

	"template/src/auth"
	"template/src/domain"

	"github.com/danielgtaylor/huma/v2"
)

func Domain(api huma.API, userStore domain.UserStore, store domain.Store, authHandler auth.AuthHandler, logger *slog.Logger) {
	register(
		api,
		makeSecuredOperation(http.MethodGet, "/user", auth.ScopeDomain),
		func(ctx context.Context, input *struct{},
		) (*struct{ Body domain.User }, error) {
			user, err := domain.GetCTX[domain.User](ctx)
			if err != nil {
				return nil, err
			}

			return &struct{ Body domain.User }{Body: user}, nil
		},
		logger,
	)
}
