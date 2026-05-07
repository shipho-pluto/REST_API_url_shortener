package save

import (
	"log/slog"
	"net/http"
	resp "url-shotener/internal/lib/api/response"
	"url-shotener/internal/lib/logger/sl"
	"url-shotener/internal/lib/random"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 6

//go:generate go run github.com/vektra/mockery/@latest --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) error
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Handlers.URL.Save"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if msg := "failed to decode request"; err != nil {
			log.Error(msg, sl.Err(err))

			render.JSON(w, r, resp.Error(msg))
			return
		}
		log.Info("successfully decoded request", slog.Any("request", req))

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandiomString(aliasLength)
		}

		log.Info("successfully validated request")

		err = urlSaver.SaveURL(req.URL, req.Alias)
		if msg := "failed to complite request"; err != nil {
			log.Error(msg, sl.Err(err))

			render.JSON(w, r, resp.Error(msg))
			return
		}
		log.Info("successfully complited request")

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
