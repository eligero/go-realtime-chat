package main

import (
	"log"
	"net/http"

	"github.com/eligero/go-realtime-chat/trace"
	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

const (
	socketBubberSize  = 1024
	messageBufferSize = 256
)

// In order to use web sockets
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBubberSize,
	WriteBufferSize: socketBubberSize,
}

type room struct {
	// channel that holds incoming messages that should be forwarded to the other clients
	forward chan *message

	// channel for clients wishing to join this room. Safety add and remove clients from clients map
	join chan *client

	// channel for clients wishing to leave this room. Safety add and remove clients from clients map
	leave chan *client

	// hold all current clients in this room
	clients map[*client]bool

	// tracer will receive trace information of activity in the room
	tracer trace.Tracer

	// avatar is how Avatar information will be obtained
	avatar Avatar
}

// newRoom makes a new room, creating the channels and map needed to be created
func newRoom(avatar Avatar) *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
		avatar:  avatar,
	}
}

func (r *room) run() {
	for {
		select { // r.clients map is only ever modified by one thing at a time
		case client := <-r.join: // joining
			r.clients[client] = true
			r.tracer.Trace("New client joined")
		case client := <-r.leave: // leaving
			delete(r.clients, client) // delete client key from clients map
			close(client.send)        // close the channel
			r.tracer.Trace("Client left")
		case msg := <-r.forward: // forward message to all clients
			r.tracer.Trace("Message received: ", msg.Message)
			for client := range r.clients {
				// Add the message to each client's send channel
				client.send <- msg
				r.tracer.Trace("-- sent to client")
			}
		}
	}
}

// room acts as a handler
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// get the socket
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}

	// create the client
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}

	// pass the client into the join channel for the current room
	r.join <- client

	// defer leaving operation for when the client is finished
	defer func() { r.leave <- client }()

	// goroutine
	go client.write()

	// Block operations, but keeping the connection alive, until is time to close it
	client.read()
}
