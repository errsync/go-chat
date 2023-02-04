package main

import (
	"errors"
	"io/ioutil"
	"path"
)

// ErrNoAvatarURL to błąd, zwracany gdy instancja Avatar
// nie pozwala na określenie i zwrócoenie adresu URL awatara.
var ErrNoAvatarURL = errors.New("chat: Nie można pobrać adresu URL awatara.")

// Avatar reprezentuje typ pozwalający na reprezentowanie
// zdjęć profilowych użytkowników.
type Avatar interface {
  // Metoda GetAvatarURL pobiera adres URL awatara dla konkretnego
  // klienta, bądź też, jeśli coś poszło źle, zwraca błąd.
  // Błąd ErrNoAvatarURL jest zwracany jeśli obiekt nie pozwala 
  // na zwrócenie adresu URL awatara dla danego klienta.
  GetAvatarURL(ChatUser) (string, error)
}

type TryAvatars []Avatar

func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
  for _, avatar := range a {
    if url, err := avatar.GetAvatarURL(u); err == nil {
      return url, nil
    }
  }
  return "", ErrNoAvatarURL
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
  if files, err := ioutil.ReadDir("avatars"); err == nil {
    for _, file := range files {
      if file.IsDir() {
        continue
      }
      if match, _ := path.Match(u.UniqueID()+"*", file.Name());
      match {
        return "/avatars/" + file.Name(), nil
      }
    }
  }
  return "", ErrNoAvatarURL
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
  url := u.AvatarURL()
  if len(url) == 0 {
    return "", ErrNoAvatarURL
  }
  return url, nil
}
// GravatarAvatar 
type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
  return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}
 

