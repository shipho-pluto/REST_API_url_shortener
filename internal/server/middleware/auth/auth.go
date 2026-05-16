package auth

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/jwt"
	"url-shortener/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

var publicPaths = []string{
	"/url/login/",
	"/url/register/",
}

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/auth"),
		)

		log.Info("auth middleware enable")

		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			if !requiresAuth(r.URL.Path) {
				next.ServeHTTP(ww, r)
				return
			}

			token := catchTocken(r)

			if token == "" {
				log.Warn("token not found in request")
				render.JSON(w, r, resp.Error("authorization required"))
				return
			}

			err := jwt.ValidToken(token)

			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					log.Info("token expired", sl.Err(jwt.ErrTokenExpired))

					render.JSON(w, r, resp.Error("relogin again your session ended"))
					return
				}
				log.Warn("ivalid token", sl.Err(err))

				render.JSON(w, r, resp.Error("authorization required"))
				return
			}

			log.Info("token is valid, client successfully authorized")
			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}

func catchTocken(r *http.Request) string {
	if cookie, err := r.Cookie("token"); err == nil {
		return cookie.Value
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		token = strings.TrimPrefix(token, "bearer ")
		if token != authHeader {
			return token
		}
	}

	return ""
}

func requiresAuth(path string) bool {
	for _, publicPath := range publicPaths {
		if path == publicPath {
			return false
		}
	}
	return strings.HasPrefix(path, "/url")
}
