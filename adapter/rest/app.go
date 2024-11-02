package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/todennus/shared/config"
	"github.com/todennus/shared/middleware"
	"github.com/todennus/user-service/wiring"
)

func App(
	config *config.Config,
	usecases *wiring.Usecases,
) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.SetupContext(config))
	r.Use(middleware.Recoverer())
	r.Use(middleware.LogRequest(config))
	r.Use(middleware.Timeout(config))
	r.Use(middleware.Authentication(config.TokenEngine))
	r.Use(middleware.WithSession(config.SessionManager))

	r.Route("/users", NewUserAdapter(usecases.UserUsecase, usecases.AvatarUsecase).Router)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })

	return r
}
