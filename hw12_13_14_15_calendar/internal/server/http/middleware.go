package internalhttp

import (
	"log"
	"net/http"
	"time"
)

type statusCodeCaptor struct {
	http.ResponseWriter
	statusCode int
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writerWrap := &statusCodeCaptor{ResponseWriter: w}
		start := time.Now()
		next.ServeHTTP(writerWrap, r)
		latency := time.Since(start)
		log.Printf(
			//nolint:lll
			"IP: %s, Дата и время запроса: %s, Метод: %s, Path: %s, HTTP Version: %s, Код ответа: %d, Latency: %v, User Agent: %s\n",
			r.RemoteAddr,
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			r.Proto,
			writerWrap.statusCode,
			latency,
			r.UserAgent(),
		)
	})
}

func (captor *statusCodeCaptor) WriteHeader(code int) {
	captor.statusCode = code
	captor.ResponseWriter.WriteHeader(code)
}
