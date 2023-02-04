package main

import (
  "time"

  "github.com/gorilla/websocket"
)

// Typ client reprezentuje pojedynczego użytkownika 
// prowadzącego konwersację z użyciem komunikatora.
type client struct {

  // socket to gniazdo internetowe do obsługi danego klienta.
  socket *websocket.Conn

  // send to kanał którym są przesyłane komunikaty.
  send chan *message

  // room to pokój rozmów używany przez klienta.
  room *room

  // pole userData zawiera informacje o użytkwoniku
  userData map[string]interface{}  
}

func (c *client) read() {
  defer c.socket.Close()
  for {
    var msg *message
    err := c.socket.ReadJSON(&msg)
    if err != nil {
      return
    }
    msg.When = time.Now()
    msg.Name = c.userData["name"].(string)
    if avatarURL, ok := c.userData["avatar_url"]; ok {
      msg.AvatarURL = avatarURL.(string)
    }
    c.room.forward <- msg
  }
}

func (c *client) write() {
  defer c.socket.Close()
  for msg := range c.send {
    err := c.socket.WriteJSON(msg)
    if err != nil {
      return
    }
  }
}
