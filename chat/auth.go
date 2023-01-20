package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		// Not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		// Error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// auth success, call the next handler
	h.next.ServeHTTP(w, r)
}

// Helper. Creates autHandler that wraps any other http.Handler
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler handles third party login process
// format: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	action := segments[2]
	provider := segments[3]

	// Improvement: Manage panic if segments not follow the format

	switch action {
	case "login":
		log.Println("TODO: handle login for ", provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}
