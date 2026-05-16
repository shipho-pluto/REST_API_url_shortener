package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enable")

		fn := func(w http.ResponseWriter, r *http.Request) {
			// Начальная инфа
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_address", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			// Куда будут записываться следующие запросы по цепочке
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			timer := time.Now()
			// Под конец когда мы вернулись к этому middleware мы выписываем инфу
			defer func() {
				entry.Info("request complited",
					slog.Int("sratus", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(timer).String()),
				)
			}()

			// Передача управления следующиму middleware
			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
