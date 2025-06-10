package middleware

import (
	"log/slog"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

func CreateErrorInterceptor(api huma.API, logger *slog.Logger) huma.API {
	return &errorInterceptorAPI{
		API:    api,
		logger: logger,
	}
}

type errorInterceptorAPI struct {
	huma.API
	logger *slog.Logger
}

func (a *errorInterceptorAPI) Transform(ctx huma.Context, status string, response any) (any, error) {
	if len(status) > 0 && status[0] >= '4' && status[0] < '5' {
		if response != nil {
			a.logErrorDetails(status, response)
		}
	}

	return a.API.Transform(ctx, status, response)
}

func (a *errorInterceptorAPI) logErrorDetails(status string, response any) {
	rv := reflect.ValueOf(response)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.Struct {
		errorsField := rv.FieldByName("Errors")
		if errorsField.IsValid() && errorsField.Kind() == reflect.Slice && errorsField.Len() > 0 {
			a.logger.Error("Validation errors", "status", status, "errors", errorsField.Interface())
			return
		}

		detailsField := rv.FieldByName("Details")
		if detailsField.IsValid() {
			a.logger.Error("Validation errors", "status", status, "details", detailsField.Interface())
			return
		}

		msgField := rv.FieldByName("Message")
		if msgField.IsValid() && msgField.Kind() == reflect.String {
			a.logger.Error("Validation error", "status", status, "message", msgField.String())
			return
		}
	}

	a.logger.Error("Validation error (status %s): %v", status, response)
}
