package hub

import (
	"sync/atomic"

	"github.com/YutakaHorikawa/gows/ws"
)

type HubManager struct {
	hubs []*Hub
}

type Hub struct {
	clients        map[string]map[*ws.Client]bool
	Broadcast      chan *ws.Client
	Register       chan *ws.Client
	Unregister     chan *ws.Client
	connectedAmout int32
}

func (hm *HubManager) GetHub() *Hub {
	var h *Hub
	first := true
	for i := range hm.hubs {
		if first == true {
			h = hm.hubs[i]
			first = false
		} else {
			if hm.hubs[i].connectedAmout < h.connectedAmout {
				h = hm.hubs[i]
			}
		}
	}

	return h
}

func (hm *HubManager) GetHubByRoomid(roomId string) *Hub {
	var h *Hub
	for i := range hm.hubs {
		hub := hm.hubs[i]
		if _, ok := hub.clients[roomId]; ok {
			return hub
		}
	}

	return h
}

func (hm *HubManager) setHub(hub *Hub, index int) {
	hm.hubs[index] = hub
}

func (hm *HubManager) RunAllHub() {
	for i := range hm.hubs {
		go hm.hubs[i].run()
	}
}

func (h *Hub) IncreaseConnectedAmount() {
	atomic.AddInt32(&h.connectedAmout, 1)
}

func NewHubManager(worker int) *HubManager {
	hubManager := &HubManager{
		hubs: make([]*Hub, worker),
	}

	for i := 0; i < worker; i++ {
		hub := newHub()
		hubManager.setHub(hub, i)
	}

	return hubManager
}

func newHub() *Hub {
	return &Hub{
		Broadcast:      make(chan *ws.Client),
		Register:       make(chan *ws.Client),
		Unregister:     make(chan *ws.Client),
		clients:        make(map[string]map[*ws.Client]bool),
		connectedAmout: 0,
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.Register:
			if _, ok := h.clients[client.roomId]; ok {
				h.clients[client.roomId][client] = true
			} else {
				h.clients[client.roomId] = make(map[*ws.Client]bool)
				h.clients[client.roomId][client] = true
			}
			h.IncreaseConnectedAmount()
		case client := <-h.Unregister:
			if _, ok := h.clients[client.roomId][client]; ok {
				delete(h.clients[client.roomId], client)
				close(client.send)
			}
		case client := <-h.Broadcast:
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
