package registerroute

import (
	"context"
	"log/slog"

	"template/src/auth"

	"github.com/danielgtaylor/huma/v2"
)

func makeOperation(method string, path string) huma.Operation {
	op := huma.Operation{
		Method: method,
		Path:   path,
		Tags:   []string{"public"},
	}

	return op
}

func makeSecuredOperation(method string, path string, scope auth.Scope) huma.Operation {
	op := makeOperation(method, path)
	op.Tags = []string{string(scope)}
	op.Security = []map[string][]string{
		{
			"scope": {string(scope)},
		},
	}

	return op
}

func register[T any, I any](api huma.API, op huma.Operation, handler func(context.Context, *I) (*T, error), logger *slog.Logger) {
	operationString := op.Method + " " + op.Path

	huma.Register(api, op, func(ctx context.Context, input *I) (*T, error) {
		output, err := handler(ctx, input)
		if err != nil {
			logger.Error("Call failed", "operation", operationString, "error", err)
			return nil, err
		}

		logger.Info("Call succeeded", "operation", operationString)

		return output, err
	})
}
