package main

import (
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func uploaderHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userid")
	file, header, err := r.FormFile("avatarFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = os.Stat("avatars")
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir("avatars", 0755)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	// Write file
	filename := path.Join("avatars", userId+path.Ext(header.Filename))
	err = os.WriteFile(filename, data, 0777)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Remove previous uploaded avatars
	files, err := filepath.Glob("avatars/" + userId + "*")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, f := range files {
		if f != filename {
			if err := os.Remove(f); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	io.WriteString(w, "Successful")
}
