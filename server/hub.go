package main

type Hub struct {
	clients    map[string]map[*Client]bool
	broadcast  chan *Client
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			if _, ok := h.clients[client.roomId]; ok {
				h.clients[client.roomId][client] = true
			} else {
				h.clients[client.roomId] = make(map[*Client]bool)
				h.clients[client.roomId][client] = true
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client.roomId][client]; ok {
				delete(h.clients[client.roomId], client)
				close(client.send)
			}
		case client := <-h.broadcast:
			message := client.message
			for c := range h.clients[client.roomId] {
				select {
				case c.send <- message:
				default:
					close(c.send)
					delete(h.clients[c.roomId], c)
				}
			}
		}
	}
}
