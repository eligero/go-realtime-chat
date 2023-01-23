package main

import (
	"net/http"
	"strings"
)

// imageHandler handles profile image
// format: /imager/action
func imageHandler(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	action := segments[2]

	switch action {
	case "file":
		avatars = UseFileSystemAvatar
	case "gravatar":
		avatars = UseGravatar
	case "provider":
		avatars = UseAuthAvatar
	}

	w.Header().Set("Location", "/logout")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
