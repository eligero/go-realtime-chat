package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once       sync.Once
	filename   string
	thTemplate *template.Template
}

// templateHandler method ServeHTTP handles the HTTP request. It satisfies the http.Handler interface
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// sync.Once type guarantees that the function we pass as an argument will only be executed once, regardless of how many goroutines are calling ServeHTTP
	t.once.Do(func() { // load, compile and execute the template
		t.thTemplate = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.thTemplate.Execute(w, nil) // Write the output to the ResponseWriter
}

func main() {
	http.Handle("/", &templateHandler{filename: "chat.html"})

	// Start the web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
