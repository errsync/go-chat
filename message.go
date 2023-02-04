package main

import (
  "time"
)

// reprezentuje pojedynczą wiadomość
type message struct {
  Name      string
  Message   string
  When      time.Time
  AvatarURL string
}
