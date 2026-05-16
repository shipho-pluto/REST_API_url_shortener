package login

import (
	"context"
	"log/slog"
	"net/http"
	rest "url-shortener/internal/lib/api/request"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"

	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
	Token string `json:"token"`
}

type Loginer interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int32,
	) (token string, err error)
	AppID() int32
}

func New(ctx context.Context, log *slog.Logger, cls Loginer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Handlers.URL.Login"

		log = log.With(
			slog.String("op", op),
		)

		var req rest.Request
		err := render.DecodeJSON(r.Body, &req)
		if msg := "failed to decode request"; err != nil {
			log.Error(msg, sl.Err(err))

			render.JSON(w, r, resp.Error("request required"))
			return
		}

		log.Info("successfully got credentionals",
			slog.Any("email", req.Email),
			slog.Any("password", req.Password),
		)

		token, err := cls.Login(ctx, req.Email, req.Password, cls.AppID())
		if err != nil {
			log.Error("error with sso client", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))
			return
		}
		log.Info("successfully complited request")

		render.JSON(w, r, respOK(token))
	}
}

func respOK(token string) *Response {
	return &Response{
		Response: resp.OK(),
		Token:    token,
	}
}
