package middleware

import (
	"net/http"
	"time"
)

type Logger interface {
	Log(format string, info ...any)
}

type LoggerMiddleware struct {
	logger Logger
}

func NewLoggerMiddleware(logger Logger) LoggerMiddleware {
	return LoggerMiddleware{logger: logger}
}
func (lm LoggerMiddleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		lm.logger.Log("%s %s %s", req.Method, req.RequestURI, time.Since(start))
	})
}
