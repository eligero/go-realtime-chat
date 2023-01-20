package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
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
	t.thTemplate.Execute(w, r) // Write the output to the ResponseWriter
}

func main() {
	var addr = flag.String("addr", ":8080", "The address of the application")
	flag.Parse() // parse the flags

	gomniauth.SetSecurityKey(os.Getenv("security_key"))
	gomniauth.WithProviders(
		google.New(os.Getenv("google_key"), os.Getenv("google_secret"), "http://localhost:8080/auth/callback/google"),
		github.New(os.Getenv("github_key"), os.Getenv("github_secret"), "http://localhost:8080/auth/callback/github"),
	)

	r := newRoom()

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	// Run the room in a goroutine
	go r.run()

	// Start the web server
	log.Println("Starting the web server on ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
