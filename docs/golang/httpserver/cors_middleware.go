package httpserver

import (
	"net/http"
	"regexp"
	"strings"

	"gitlab.nordstrom.com/sentry/gohttp/middleware"
)

type corsMiddleware struct {
	regex string
}

//NewCorsMiddleware returns a new cors middleware
func NewCorsMiddleware(regex string) middleware.Middleware {
	return &corsMiddleware{regex}
}

func (m corsMiddleware) Wrap(h http.Handler) http.Handler {
	return newCorsHandler(h, m.regex)
}

type corsHandler struct {
	regex         *regexp.Regexp
	innerHandler  http.Handler
	serveHTTPFunc func(http.ResponseWriter, *http.Request)
}

//NewCorsHandler returns a new cors handler with the provided inner handler
func newCorsHandler(innerHandler http.Handler, regexString string) corsHandler {
	corsHandler := corsHandler{nil, innerHandler, nil}
	if strings.TrimSpace(regexString) == "" {
		return corsHandler
	}
	corsHandler.regex = regexp.MustCompile(regexString)
	corsHandler.serveHTTPFunc = func(w http.ResponseWriter, r *http.Request) {
		if corsHandler.regex != nil {
			o := r.Header.Get("Origin")
			if corsHandler.regex.FindStringIndex(o) != nil {
				w.Header().Add("Access-Control-Allow-Origin", o)
			}
		}
		corsHandler.innerHandler.ServeHTTP(w, r)
	}
	return corsHandler
}

func (h corsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.serveHTTPFunc != nil {
		h.serveHTTPFunc(w, r)
	}
}

func NewCorsOptionsHandler(regexString string) corsHandler {
	corsHandler := corsHandler{nil, nil, nil}
	if strings.TrimSpace(regexString) == "" {
		return corsHandler
	}
	corsHandler.regex = regexp.MustCompile(regexString)
	corsHandler.serveHTTPFunc = func(w http.ResponseWriter, r *http.Request) {
		if corsHandler.regex != nil {
			o := r.Header.Get("Origin")
			if corsHandler.regex.FindStringIndex(o) != nil {
				w.Header().Add("Access-Control-Allow-Origin", o)
				w.Header().Add("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
				w.Header().Add("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
			}
		}
	}
	return corsHandler
}
