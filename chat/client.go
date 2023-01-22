package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// client represents a single user chatting
type client struct {
	// web socket for this client
	socket *websocket.Conn
	// buffered channel on which messages are sent
	send chan *message
	// room where client is chatting
	room *room
	// userData holds information about the user
	userData map[string]interface{}
}

// read from the socked via ReadMessage method
func (c *client) read() {
	// close the socket if an error occurs
	defer c.socket.Close()
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err != nil {
			return
		}
		msg.TimeStamp = time.Now()
		msg.Name = c.userData["name"].(string)
		if avatarUrl, ok := c.userData["avatar_url"]; ok {
			msg.AvatarURL = avatarUrl.(string)
		}
		c.room.forward <- msg
	}
}

// write accepts messages from send channel, writing everithing out of the socket via WriteMessage method
func (c *client) write() {
	// close the socket if an error occurs
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}
