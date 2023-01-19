package main

import (
	"github.com/gorilla/websocket"
)

// client represents a single user chatting
type client struct {
	// web socket for this client
	socket *websocket.Conn
	// buffered channel on which messages are sent
	send chan []byte
	// roome where client is chatting
	room *room
}

// read from the socked via ReadMessage method
func (c *client) read() {
	// close the socket if an error occurs
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

// write accepts messages from send channel, writing everithing out of the socket via WriteMessage method
func (c *client) write() {
	// close the socket if an error occurs
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
