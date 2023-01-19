package main

type room struct {
	// channel that holds incoming messages that should be forwarded to the other clients
	forward chan []byte

	// channel for clients wishing to join this room. Safety add and remove clients from clients map
	join chan *client

	// channel for clients wishing to leave this room. Safety add and remove clients from clients map
	leave chan *client

	// hold all current clients in this room
	clients map[*client]bool
}

func (r *room) run() {
	for {
		select { // r.clients map is only ever modified by one thing at a time
		case client := <-r.join: // joining
			r.clients[client] = true
		case client := <-r.leave: // leaving
			delete(r.clients, client) // delete client key from clients map
			close(client.send)        // close the channel
		case msg := <-r.forward: // forward message to all clients
			for client := range r.clients {
				// Add the message to each client's send channel
				client.send <- msg
			}
		}
	}
}
