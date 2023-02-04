package main

import (
  "io/ioutil"
  "os"
  "path"

  gomniauthtest "github.com/stretchr/gomniauth/test"

  "testing"
)

func TestAuthAvatar(t *testing.T) {

  var authAvatar AuthAvatar
  testUser := &gomniauthtest.TestUser{}
  testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
  testChatUser := &chatUser{User: testUser}
  url, err := authAvatar.GetAvatarURL(testChatUser)
  if err != ErrNoAvatarURL {
    t.Error("Gdy nie jest dostępna żadna wartość, metoda AuthAvatar.GetAvatarURL powinna zwrócić ErrNoAvatarURL!")
  }

  testURL := "http://adresURL-usługi-gravatar/"
  testUser = &gomniauthtest.TestUser{}
  testChatUser.User = testUser
  testUser.On("AvatarURL").Return(testURL, nil)
  url, err = authAvatar.GetAvatarURL(testChatUser)
  if err != nil {
    t.Error("Metoda AuthAvatar.GetAvatarURL nie powinna zwracać błędów, gdy wartość jest określona!")
  }
  if url != testURL {
    t.Error("Metoda AuthAvatar.GetAvatarURL powinna zwrócić prawidłowy adres URL")
  }
}


func TestGravatarAvatar(t *testing.T) {

  var gravatarAvatar GravatarAvatar
  user := &chatUser{uniqueID: "abc"}

  url, err := gravatarAvatar.GetAvatarURL(user)
  if err != nil {
    t.Error("Metoda GravatarAvatar.GetAvatarURL nie powinna zwrócić błędu!")
  }
  if url != "//www.gravatar.com/avatar/abc" {
    t.Errorf("Metoda GravatarAvatar.GetAvatarURL błędnie zwróciła %s", url)
  }
}


func TestFileSystemAvatar(t *testing.T) {

  // utworzenie testowego pliku graficznego
  filename := path.Join("avatars", "abc.jpg")
  if err := os.MkdirAll("avatars", 0777); err != nil {
    t.Errorf("Nie udało się utworzyć katalogu 'avatars': %s", err)
  }
  if err := ioutil.WriteFile(filename, []byte{}, 0777); err != nil {
    t.Errorf("Nie udało się utworzyć pliku awatara: %s", err)
  }
  defer os.Remove(filename)

  var fileSystemAvatar FileSystemAvatar
  user := &chatUser{uniqueID: "abc"}

  url, err := fileSystemAvatar.GetAvatarURL(user)
  if err != nil {
    t.Errorf("Funkcja FileSystemAvatar.GetAvatarURL nie powinna zwrócić błędu: %s", err)
  }
  if url != "/avatars/abc.jpg" {
    t.Errorf("Funkcja FileSystemAvatar.GetAvatarURL błędnie zwróciła %s", url)
  }
}
