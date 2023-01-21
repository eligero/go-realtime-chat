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
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
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

	data := map[string]interface{}{
		"Host": r.Host,
	}

	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.thTemplate.Execute(w, data) // Write the output to the ResponseWriter
}

func main() {
	var addr = flag.String("addr", ":8080", "The address of the application")
	flag.Parse() // parse the flags

	gomniauth.SetSecurityKey(os.Getenv("security_key"))
	gomniauth.WithProviders(
		google.New(os.Getenv("google_key"), os.Getenv("google_secret"), "http://localhost:8080/auth/callback/google"),
		github.New(os.Getenv("github_key"), os.Getenv("github_secret"), "http://localhost:8080/auth/callback/github"),
		facebook.New(os.Getenv("facebook_key"), os.Getenv("facebook_secret"), "http://localhost:8080/auth/callback/facebook"),
	)

	r := newRoom()

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "", // Clear Cookie's value
			Path:   "/",
			MaxAge: -1, // Immediately removed by the browser, if supported by browser
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	// Run the room in a goroutine
	go r.run()

	// Start the web server
	log.Println("Starting the web server on ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
