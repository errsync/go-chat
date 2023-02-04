package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/errsync/go-chat/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"

  "github.com/errsync/go-chat/authkey"
)

// zmienna określająca aktywną implementację Avatar
var avatars Avatar = TryAvatars{
  UseFileSystemAvatar, 
  UseAuthAvatar,
  UseGravatar}

// struktura reprezentująca pojedyczny szablon.
type templateHandler struct {
  filename string
  templ    *template.Template 
}

// ServeHTTP obsługuje żądania HTTP.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if t.templ == nil {
    t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
  }

  data := map[string]interface{}{
    "Host": r.Host,
  }
  if authCookie, err := r.Cookie("auth"); err == nil {
    data["UserData"] = objx.MustFromBase64(authCookie.Value)
  }

  t.templ.Execute(w, data)
}

var addr = flag.String("host", ":8080", "Komputer na którym działa aplikacja.")

func main() {
  
  flag.Parse() // analiza flag wiersza poleceń

  // konfiguracja pakietu gomniauth
  gomniauth.SetSecurityKey("98dfbg7iu2hdbf56evihjw4tuiyub34noilk")
  gomniauth.WithProviders(
    github.New(authkey.FbID, authkey.GithubKEY, "http://localhost:8080/auth/callback/github"),
    google.New(authkey.GoogleID, authkey.GoogleKEY, "http://localhost:8080/auth/callback/google"),
    facebook.New(authkey.FbID, authkey.FbKEY, "http://localhost:8080/auth/callback/facebook"),
  )

  r := newRoom()
  r.tracer = trace.New(os.Stdout)

  http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
  http.Handle("/login", &templateHandler{filename: "login.html"})
  http.HandleFunc("/auth/", loginHandler)
  http.Handle("/room", r)
  http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
    http.SetCookie(w, &http.Cookie{
      Name:   "auth",
      Value:  "",
      Path:   "/",
      MaxAge: -1,
    })
    w.Header().Set("Location", "/chat")
    w.WriteHeader(http.StatusTemporaryRedirect)
  })  

  http.Handle("/upload", &templateHandler{filename: "upload.html"})
  http.HandleFunc("/uploader", uploaderHandler)

  http.Handle("/avatars/",
    http.StripPrefix("/avatars/",
      http.FileServer(http.Dir("./avatars"))))

  // uruchomienie pokoju rozmów
  go r.run()

  // uruchomienie serwera WWW
  log.Println("Uruchamianie serwera WWW pod adresem", *addr)
  if err := http.ListenAndServe(*addr, nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }

}
