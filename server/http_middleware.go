package server

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

type HTTPMiddleware func(http.Handler) http.Handler

type LoggerHandler struct {
	Logger *logrus.Logger
	Next   http.Handler
}

type LoggerResponseWriter struct {
	http.ResponseWriter
	Code int
}

func (w *LoggerResponseWriter) WriteHeader(status int) {
	w.Code = status
	w.ResponseWriter.WriteHeader(status)
}

// 日志格式
// 2017-01-01 10:00:00 / 10.10.10.10 [200] 10ms
func (h *LoggerHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	lw := &LoggerResponseWriter{w, http.StatusOK}
	defer func(b time.Time) {
		if request.Method == http.MethodPost {
			request.ParseForm()
			h.Logger.Infof("%s\t%s\t%s\t%s\t%d %s\t%v", request.Method, request.URL.Path, request.PostForm.Encode(), GetIPAdress(request), lw.Code, http.StatusText(lw.Code), time.Since(b))
		} else {
			h.Logger.Infof("%s\t%s\t%s\t%s\t%d %s\t%v", request.Method, request.URL.Path, request.URL.Query().Encode(), GetIPAdress(request), lw.Code, http.StatusText(lw.Code), time.Since(b))
		}
	}(time.Now())
	h.Next.ServeHTTP(lw, request)
}

func LoggerHandlerMiddleware(l *logrus.Logger) HTTPMiddleware {
	return func(h http.Handler) http.Handler {
		return &LoggerHandler{Logger: l, Next: h}
	}
}
