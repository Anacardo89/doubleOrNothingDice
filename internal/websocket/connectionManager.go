package websocket

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Connection struct {
	Conn         *websocket.Conn
	LastActivity time.Time
}

type ConnectionManager struct {
	connections map[string]*Connection
	mu          sync.RWMutex
	timeout     time.Duration
}

func NewConnectionManager(timeoutMinutes int) *ConnectionManager {
	cm := &ConnectionManager{
		connections: make(map[string]*Connection),
		timeout:     time.Duration(timeoutMinutes) * time.Minute,
	}
	cm.startTimeoutWatcher()
	return cm
}

func (cm *ConnectionManager) Add(clientID string, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if existingConn, ok := cm.connections[clientID]; ok {
		existingConn.Conn.Close()
	}
	cm.connections[clientID] = &Connection{
		Conn:         conn,
		LastActivity: time.Now(),
	}
}

func (cm *ConnectionManager) Remove(clientID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if conn, ok := cm.connections[clientID]; ok {
		conn.Conn.Close()
		delete(cm.connections, clientID)
	}
}

func (cm *ConnectionManager) Get(clientID string) (*websocket.Conn, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	conn, ok := cm.connections[clientID]
	if !ok {
		return nil, false
	}
	return conn.Conn, true
}

func (cm *ConnectionManager) UpdateActivity(clientID string) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	if conn, ok := cm.connections[clientID]; ok {
		conn.LastActivity = time.Now()
	}
}

func (cm *ConnectionManager) startTimeoutWatcher() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			cm.mu.Lock()
			for clientID, conn := range cm.connections {
				if time.Since(conn.LastActivity) > cm.timeout {
					log.Printf("Connection timeout: closing connection for user %s", clientID)
					conn.Conn.Close()
					delete(cm.connections, clientID)
				}
			}
			cm.mu.Unlock()
		}
	}()
}
