package redirect

import (
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/@latest --name=URLProvider
type URLProvider interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlProvider URLProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Handlers.URL.Get"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if msg := "alias if empty"; alias == "" {
			log.Info(msg)

			render.JSON(w, r, resp.Error(msg))
			return
		}
		log.Info("successfully got alias", slog.Any("alias", alias))

		url, err := urlProvider.GetURL(alias)
		if msg := "failed to complete request"; err != nil {
			log.Error(msg, sl.Err(err))

			render.JSON(w, r, resp.Error(msg))
			return
		}
		log.Info("successfully got URL from provider")

		http.Redirect(w, r, url, http.StatusFound)
	}
}
