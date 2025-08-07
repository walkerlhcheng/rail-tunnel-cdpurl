package handlers

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"rail-tunnel/managers"
	"rail-tunnel/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type TunnelHandlers struct {
	connectionManager *managers.ConnectionManager
	upgrader          websocket.Upgrader
	pendingRequests   map[string]*PendingRequest
	pendingMu         sync.RWMutex
}

type PendingRequest struct {
	ResponseChan chan *models.Message
	Timeout      time.Time
}

func NewTunnelHandlers(cm *managers.ConnectionManager) *TunnelHandlers {
	return &TunnelHandlers{
		connectionManager: cm,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins
			},
		},
		pendingRequests: make(map[string]*PendingRequest),
	}
}

// HandleWebSocket handles WS /ws/connect - Auto-configure tunnel
func (h *TunnelHandlers) HandleWebSocket(c *gin.Context) {
	// Extract connection info from query params (sent by CLI)
	localPort := c.Query("port") // ?port=3000
	if localPort == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'port' parameter"})
		return
	}

	// Auto-configure connection
	serverURL := getServerURL(c.Request)
	req := models.ConnectRequest{
		LocalPort: parsePort(localPort),
		LocalURL:  fmt.Sprintf("http://localhost:%s", localPort),
	}
	connection := h.connectionManager.Connect(req, serverURL)
	log.Printf("Auto-configured tunnel: %s -> %s", connection.PublicURL, connection.LocalURL)

	// Upgrade to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// Store connection
	connection.Conn = conn
	h.connectionManager.SetConnected(true)
	log.Printf("WebSocket connected: %s -> %s", connection.PublicURL, connection.LocalURL)

	// Handle messages
	for {
		var msg models.Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		h.handleWebSocketMessage(&msg)
	}

	// Clean up connection
	h.connectionManager.SetConnected(false)
	log.Printf("WebSocket disconnected")
}

// HandleTunnelTraffic handles ALL HTTP requests as proxy
func (h *TunnelHandlers) HandleTunnelTraffic(c *gin.Context) {
	// Get the current connection
	connection, connected := h.connectionManager.GetConnection()
	if !connected {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "Tunnel not available",
			"message": "No tunnel connected. Start your CLI client first.",
		})
		return
	}

	// Forward ALL requests to tunnel client
	requestPath := c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		requestPath += "?" + c.Request.URL.RawQuery
	}

	h.forwardRequestToTunnel(c, connection, requestPath)
}

// handleWebSocketMessage handles incoming WebSocket messages
func (h *TunnelHandlers) handleWebSocketMessage(msg *models.Message) {
	switch msg.Type {
	case "ping":
		// Respond to ping with pong
		log.Printf("Received ping from tunnel")
	case "pong":
		log.Printf("Received pong from tunnel")
	case "http_response":
		// Handle HTTP response from tunnel client
		log.Printf("Received HTTP response for request %s: status=%d", msg.RequestID, msg.StatusCode)
		h.handleResponse(msg)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

// forwardRequestToTunnel forwards HTTP request to tunnel client via WebSocket
func (h *TunnelHandlers) forwardRequestToTunnel(c *gin.Context, connection *models.TunnelConnection, requestPath string) {
	requestID := uuid.New().String()

	// Read request body
	var body interface{}
	if c.Request.Body != nil {
		rawBody := make([]byte, c.Request.ContentLength)
		c.Request.Body.Read(rawBody)
		body = string(rawBody)
	}

	msg := models.Message{
		Type:      "http_request",
		RequestID: requestID,
		Method:    c.Request.Method,
		URL:       requestPath,
		Headers:   c.Request.Header,
		Body:      body,
	}

	if err := connection.Conn.WriteJSON(msg); err != nil {
		log.Printf("Failed to forward request: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Tunnel unavailable"})
		return
	}

	log.Printf("Forwarded %s %s to tunnel", c.Request.Method, requestPath)

	// Wait for response
	h.pendingMu.Lock()
	h.pendingRequests[requestID] = &PendingRequest{
		ResponseChan: make(chan *models.Message),
		Timeout:      time.Now().Add(10 * time.Second),
	}
	h.pendingMu.Unlock()

	select {
	case response := <-h.pendingRequests[requestID].ResponseChan:
		// Set headers from tunnel response
		if headers, ok := response.Headers.(map[string]interface{}); ok {
			for key, value := range headers {
				if headerValue, ok := value.(string); ok {
					c.Header(key, headerValue)
				}
			}
		}
		
		// Set status code
		c.Status(response.StatusCode)
		
		// Simple: just return whatever the body is
		// Headers already set above, they contain Content-Type
		if bodyStr, ok := response.Body.(string); ok {
			c.String(response.StatusCode, "%s", bodyStr)
		} else {
			c.JSON(response.StatusCode, response.Body)
		}
	case <-time.After(10 * time.Second):
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Tunnel timed out"})
	}

	h.pendingMu.Lock()
	delete(h.pendingRequests, requestID)
	h.pendingMu.Unlock()
}

// handleResponse handles incoming HTTP response from tunnel client
func (h *TunnelHandlers) handleResponse(msg *models.Message) {
	h.pendingMu.RLock()
	pendingRequest, exists := h.pendingRequests[msg.RequestID]
	h.pendingMu.RUnlock()

	if exists {
		pendingRequest.ResponseChan <- msg
	}
}

// Helper functions
func getServerURL(r *http.Request) string {
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}
	return scheme + "://" + r.Host
}

func generateRandomID() string {
	// Generate random ID using crypto/rand
	b := make([]byte, 6)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func parsePort(portStr string) int {
	if port, err := strconv.Atoi(portStr); err == nil {
		return port
	}
	return 3000 // default
}
