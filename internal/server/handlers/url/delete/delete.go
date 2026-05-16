package delete

import (
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

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

		routeCtx := chi.RouteContext(r.Context())

		// 3. Проверяем URLPath (это реальный путь без параметров)
		urlPath := routeCtx.RoutePath

		log.Info("DIAGNOSTIC",
			slog.String("url", r.URL.String()),
			slog.String("path", r.URL.Path),
			slog.String("route_path", urlPath), // шаблон маршрута
			slog.String("chi_param_alias", alias),
			slog.Int("params_count", len(routeCtx.URLParams.Keys)),
		)

		if msg := "alias is empty"; alias == "" {
			log.Warn(msg)

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
