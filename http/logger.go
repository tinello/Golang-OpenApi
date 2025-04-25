package http

import (
	"log"
	"net/http"
	"time"
)

type Logger struct {
	http.Handler
}

func (l Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := newLoggedResponseWriter(w)
	startTime := time.Now()

	l.Handler.ServeHTTP(response, r)

	logRequest(r, response.statusCode, time.Now().Sub(startTime))
}

func newLoggedResponseWriter(w http.ResponseWriter) *logResponseWriter {
	return &logResponseWriter{w, http.StatusOK}
}

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *logResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func logRequest(r *http.Request, statusCode int, elapsedTime time.Duration) {
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		userAgent = "user-agent-not-present"
	}

	log.Printf(
		"[HTTP] %s %d %s %s (%s) %s \"%s\"\n",
		r.Proto,
		statusCode,
		r.Method,
		r.URL.RequestURI(),
		RemoteIPAddressFor(r),
		elapsedTime,
		userAgent,
	)
}

func RemoteIPAddressFor(request *http.Request) string {
	xRealIp := request.Header.Get("X-Real-IP")
	if xRealIp != "" {
		return xRealIp
	}

	return removePort(request.RemoteAddr)
}
