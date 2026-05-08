package delete

import (
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

const aliasLength = 6

type URLTrasher interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlTrasher URLTrasher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Handlers.URL.Delete"

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

		err := urlTrasher.DeleteURL(alias)
		if msg := "failed to complite request"; err != nil {
			log.Error(msg, sl.Err(err))

			render.JSON(w, r, resp.Error(msg))
			return
		}
		log.Info("successfully complited request")

		render.JSON(w, r, resp.OK())
	}
}
