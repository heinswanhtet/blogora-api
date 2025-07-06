package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/heinswanhtet/blogora-api/configs"
	"github.com/heinswanhtet/blogora-api/utils"
)

func init() {
	if configs.Envs.GO_ENV != "local" {
		log.SetFlags(log.LstdFlags | log.LUTC)
	}
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		coloredStatusCode := utils.GetColoredStatusCode(wrapped.statusCode)
		coloredDuration := utils.GetColoredString(time.Since(start), utils.BLUE)
		coloredMethod := utils.GetColoredHttpMethod(r.Method)

		log.Println(coloredStatusCode, coloredMethod, r.URL.RequestURI(), coloredDuration)
	})
}
