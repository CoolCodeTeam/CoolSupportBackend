package middleware

import (

	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type HandlersMiddlwares struct {
	Logger   *logrus.Logger
}



func (m *HandlersMiddlwares) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		m.Logger.WithFields(logrus.Fields{
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"work_time":   time.Since(start),
		}).Info(r.URL.Path)
	})
}

func (m *HandlersMiddlwares) PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.Logger.WithFields(logrus.Fields{
					"method":      r.Method,
					"remote_addr": r.RemoteAddr,
					"panic":       err,
				}).Error(r.URL.Path)
				http.Error(w, "Internal server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
