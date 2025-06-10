package handler

import (
	"log/slog"
	"net/http"

	"template/src/auth"
	"template/src/domain"
	"template/src/handler/middleware"
	"template/src/handler/registerroute"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humamux"
	"github.com/gorilla/mux"
)

func NewServer(
	userStore domain.UserStore,
	store domain.Store,
	authHandler *auth.AuthHandler,
	logger *slog.Logger,
) http.Handler {
	router := mux.NewRouter()

	api := middleware.CreateErrorInterceptor(
		humamux.New(router, huma.DefaultConfig("My API", "1.0.0")),
		logger,
	)
	api.UseMiddleware(middleware.Auth(authHandler.PublicKey, logger, api, userStore))

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Error("Not found", "path", r.URL.Path, "method", r.Method)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	})
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Error("Method not allowed", "path", r.URL.Path, "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	})

	registerroute.Public(api, userStore, authHandler, logger)
	registerroute.Domain(api, userStore, store, authHandler, logger)

	return router
}
