package services

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/gommon/log"
	"sync"
)

type Client struct {
	UserID int64
	Conn   *websocket.Conn
	Send   chan interface{}
}

type SocketManager struct {
	Clients map[int64]*Client
	Mutex   sync.RWMutex
}

// Global instance to be used by handlers and controllers
var Hub = &SocketManager{
	Clients: make(map[int64]*Client),
}

func (m *SocketManager) Register(client *Client) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.Clients[client.UserID] = client
}

func (m *SocketManager) Unregister(userID int64) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	if client, ok := m.Clients[userID]; ok {
		close(client.Send)
		delete(m.Clients, userID)
	}
}

func (m *SocketManager) SendToUser(userID int64, data interface{}) {
	// 1. Use Read Lock to check the map without blocking other readers
	m.Mutex.RLock()
	client, ok := m.Clients[userID]
	m.Mutex.RUnlock()

	if ok {
		// 2. Use a non-blocking select to prevent the entire Hub from hanging
		// if one user's connection is extremely slow or buffered channel is full
		select {
		case client.Send <- data:
			// Message successfully sent to the user's WritePump
		default:
			log.Printf("Warning: Client %d send buffer full, dropping message", userID)
		}
	} else {
		// User is not currently connected to the WebSocket
		log.Printf("User %d is offline, skipping real-time push", userID)
	}
}
