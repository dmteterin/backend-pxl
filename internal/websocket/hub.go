package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	canvas     *PixelCanvas
	mutex      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		canvas:     NewPixelCanvas(),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true

			// Send current user list to the new client
			userList := make([]string, 0, len(h.clients))
			for c := range h.clients {
				userList = append(userList, c.ID)
			}
			clientCount := len(h.clients)
			h.mutex.Unlock()

			log.Printf("Client %s connected. Total clients: %d", client.ID, clientCount)

			// Send current canvas state to the new client
			canvasStateMessage := NewCanvasStateMessage(h.canvas.GetPixels())
			if msgBytes, err := json.Marshal(canvasStateMessage); err == nil {
				select {
				case client.Send <- msgBytes:
				default:
					close(client.Send)
				}
			}

			// Notify other clients about the new user
			joinMessage := NewUserJoinedMessage(client.ID)
			if msgBytes, err := json.Marshal(joinMessage); err == nil {
				h.broadcastToOthers(msgBytes, client)
			}

			// Send updated user list to ALL clients (including the new one)
			userListMessage := NewUserListMessage(userList)
			if msgBytes, err := json.Marshal(userListMessage); err == nil {
				h.broadcastToAll(msgBytes)
			}

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)

				// Create updated user list after removing the client
				userList := make([]string, 0, len(h.clients))
				for c := range h.clients {
					userList = append(userList, c.ID)
				}
				clientCount := len(h.clients)
				h.mutex.Unlock()

				log.Printf("Client %s disconnected. Total clients: %d", client.ID, clientCount)

				leftMessage := NewUserLeftMessage(client.ID)
				if msgBytes, err := json.Marshal(leftMessage); err == nil {
					h.broadcastToAll(msgBytes)
				}

				// Send updated user list to all remaining clients
				updatedUserListMessage := NewUserListMessage(userList)
				if msgBytes, err := json.Marshal(updatedUserListMessage); err == nil {
					h.broadcastToAll(msgBytes)
				}
			} else {
				h.mutex.Unlock()
			}

		case message := <-h.broadcast:
			h.broadcastToAll(message)
		}
	}
}

func (h *Hub) broadcastToAll(message []byte) {
	h.mutex.RLock()
	clients := make([]*Client, 0, len(h.clients))
	for client := range h.clients {
		clients = append(clients, client)
	}
	h.mutex.RUnlock()

	var failedClients []*Client
	for _, client := range clients {
		select {
		case client.Send <- message:
		default:
			failedClients = append(failedClients, client)
		}
	}

	// Clean up failed clients
	if len(failedClients) > 0 {
		h.mutex.Lock()
		for _, client := range failedClients {
			if _, exists := h.clients[client]; exists {
				delete(h.clients, client)
				close(client.Send)
			}
		}
		h.mutex.Unlock()
	}
}

func (h *Hub) broadcastToOthers(message []byte, sender *Client) {
	h.mutex.RLock()
	clients := make([]*Client, 0, len(h.clients))
	for client := range h.clients {
		if client != sender {
			clients = append(clients, client)
		}
	}
	h.mutex.RUnlock()

	var failedClients []*Client
	for _, client := range clients {
		select {
		case client.Send <- message:
		default:
			failedClients = append(failedClients, client)
		}
	}

	// Clean up failed clients
	if len(failedClients) > 0 {
		h.mutex.Lock()
		for _, client := range failedClients {
			if _, exists := h.clients[client]; exists {
				delete(h.clients, client)
				close(client.Send)
			}
		}
		h.mutex.Unlock()
	}
}

func (h *Hub) GetClientCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}

func (h *Hub) Broadcast() chan []byte {
	return h.broadcast
}

func (h *Hub) Register() chan *Client {
	return h.register
}

func (h *Hub) Unregister() chan *Client {
	return h.unregister
}

func (h *Hub) HandlePixelMessage(msg *Message, client *Client) {
	switch msg.Type {
	case MessageTypePixelUpdate:
		if h.canvas.UpdatePixel(msg.X, msg.Y, msg.Color) {
			pixelMessage := NewPixelUpdateMessage(msg.X, msg.Y, msg.Color, client.ID)
			if msgBytes, err := json.Marshal(pixelMessage); err == nil {
				h.broadcastToAll(msgBytes)
			}
		}
	case MessageTypeCanvasClear:
		h.canvas.Clear()
		clearMessage := NewCanvasClearMessage()
		if msgBytes, err := json.Marshal(clearMessage); err == nil {
			h.broadcastToAll(msgBytes)
		}
	}
}
