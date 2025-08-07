package models

import (
	"time"

	"github.com/gorilla/websocket"
)

// TunnelConnection represents the single tunnel connection
type TunnelConnection struct {
	LocalURL  string    `json:"localUrl"`
	LocalPort int       `json:"localPort"`
	PublicURL string    `json:"publicUrl"`
	CreatedAt time.Time `json:"createdAt"`
	Conn      *websocket.Conn `json:"-"`
	Connected bool      `json:"connected"`
}

// ConnectRequest represents the request to connect tunnel
type ConnectRequest struct {
	LocalPort int    `json:"localPort"`
	LocalURL  string `json:"localUrl"`
}

// Message represents WebSocket message
type Message struct {
	Type       string      `json:"type"`
	RequestID  string      `json:"requestId,omitempty"`
	Method     string      `json:"method,omitempty"`
	URL        string      `json:"url,omitempty"`
	Headers    interface{} `json:"headers,omitempty"`
	Body       interface{} `json:"body,omitempty"`
	StatusCode int         `json:"statusCode,omitempty"`
	Timestamp  int64       `json:"timestamp,omitempty"`
}
