package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"

	gomniauthcommon "github.com/stretchr/gomniauth/common"
)

// ChatUser to interfejs uzytkownka z ID i avatarem
type ChatUser interface {
  UniqueID() string
  AvatarURL() string
}
type chatUser struct {
  gomniauthcommon.User
  uniqueID string
}

func (u chatUser) UniqueID() string {
  return u.uniqueID
}

type authHandler struct {
  next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  _, err := r.Cookie("auth")
  if err == http.ErrNoCookie {
    // brak uwierzytelnienia
    w.Header().Set("Location", "/login")
    w.WriteHeader(http.StatusTemporaryRedirect)
    return 
  }
  if err != nil {
    // jakiś inny błąd
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  // udało się - wywołujemy następny obiekt obsługi
  h.next.ServeHTTP(w, r)
}

// Metoda MustAuth pobiera obiekt obsługi, aby wymusić wykonanie uwierzytelaniania
func MustAuth(handler http.Handler) http.Handler {
  return &authHandler{next: handler}
}

// Metoda loginHandler obsługuje zewnętrzny proces logowania.
func loginHandler(w http.ResponseWriter, r *http.Request) {
  segs := strings.Split(r.URL.Path, "/")
  action := segs[2]
  provider := segs[3]
  switch action {
  case "login":

    provider, err := gomniauth.Provider(provider)
    if err != nil {
      http.Error(w, fmt.Sprintf("Błąd podczas próby pobrania dostawcy %s: %s", provider, err), http.StatusBadRequest)
      return
    }

    loginURL, err := provider.GetBeginAuthURL(nil, nil)
    if err != nil {
      http.Error(w, fmt.Sprintf("Błąd wywołania GetBeginAuthURL dla %s: %s", provider, err), http.StatusInternalServerError)
      return
    }

    w.Header().Set("Location", loginURL)
    w.WriteHeader(http.StatusTemporaryRedirect)

  case "callback":

    provider, err := gomniauth.Provider(provider)
    if err != nil {
      http.Error(w, fmt.Sprintf("Bład podczas próby pobrania dostawcy %s: %s", provider, err), http.StatusBadRequest)
      return
    }

    // pobranie danych uwierzytelniających
    creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
    if err != nil {
      http.Error(w, fmt.Sprintf("Błąd próby dokonania uwierzytelniania dla %s: %s", provider, err), http.StatusInternalServerError)
      return
    }

    // pobranie użytkownika
    user, err := provider.GetUser(creds)
    if err != nil {
      log.Fatalln("Błąd podczas próby pobrania użytkownika ", provider, "-", err)
    }
    chatUser := &chatUser{User: user}
          
    m := md5.New()
    io.WriteString(m, strings.ToLower(user.Email()))
    chatUser.uniqueID = fmt.Sprintf("%x", m.Sum(nil))
         
    avatarURL, err := avatars.GetAvatarURL(chatUser)
    if err != nil {
      log.Fatalln("Błąd podczas próby wywołania metody GetAvatarURL", "-", err)
    }
    
    authCookieValue := objx.New(map[string]interface{}{
      "userid":     chatUser.uniqueID,
      "name":       user.Name(),
      "avatar_url": avatarURL,
    }).MustBase64()

    http.SetCookie(w, &http.Cookie{
      Name:  "auth",
      Value: authCookieValue,
      Path:  "/"})

    w.Header().Set("Location", "/chat")
    w.WriteHeader(http.StatusTemporaryRedirect)

  default:
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "Akcja autoryzacyjna %s nie jest obsługiwana", action)
  }
}

