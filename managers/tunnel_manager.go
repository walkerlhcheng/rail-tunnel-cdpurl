package managers

import (
	"log"
	"sync"
	"time"

	"rail-tunnel/models"
)

// ConnectionManager manages the single tunnel connection
type ConnectionManager struct {
	connection *models.TunnelConnection
	mu         sync.RWMutex
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{}
}

// Connect establishes tunnel connection
func (cm *ConnectionManager) Connect(req models.ConnectRequest, serverURL string) *models.TunnelConnection {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Close existing connection if any
	if cm.connection != nil && cm.connection.Conn != nil {
		cm.connection.Conn.Close()
	}

	cm.connection = &models.TunnelConnection{
		LocalURL:  req.LocalURL,
		LocalPort: req.LocalPort,
		PublicURL: serverURL, // Direct server URL
		CreatedAt: time.Now(),
		Connected: false,
	}

	log.Printf("Tunnel configured: %s -> %s", serverURL, req.LocalURL)
	return cm.connection
}

// GetConnection gets the current connection
func (cm *ConnectionManager) GetConnection() (*models.TunnelConnection, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.connection, cm.connection != nil && cm.connection.Connected
}

// SetConnected marks connection as connected/disconnected
func (cm *ConnectionManager) SetConnected(connected bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.connection != nil {
		cm.connection.Connected = connected
		log.Printf("Tunnel connection status: %v", connected)
	}
}

// Disconnect closes the tunnel connection
func (cm *ConnectionManager) Disconnect() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.connection != nil {
		if cm.connection.Conn != nil {
			cm.connection.Conn.Close()
		}
		cm.connection.Connected = false
		log.Printf("Tunnel disconnected")
	}
}
