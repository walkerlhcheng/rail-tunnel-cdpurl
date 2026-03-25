/*package main

import (
	"log"
	"net/http"
	"os"

	"rail-tunnel/handlers"
	"rail-tunnel/managers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set gin mode based on environment
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize dependencies
	connectionManager := managers.NewConnectionManager()
	tunnelHandlers := handlers.NewTunnelHandlers(connectionManager)

	// Initialize gin router
	r := gin.Default()

	// System routes with unique prefixes to avoid conflicts
	r.GET("/_tunnel/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "rail-tunnel",
			"message": "Rail Tunnel server is running",
		})
	})

	// Admin/debug endpoint with unique path
	r.GET("/_tunnel/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":     "rail-tunnel",
			"description": "Ngrok-like tunnel service for Railway",
			"version":     "1.0.0",
			"endpoints": []string{
				"GET /_tunnel/health - Health check",
				"GET /_tunnel/info - Service information (this page)",
				"WS /_tunnel/ws/connect?port=3000 - WebSocket connection for tunnel client",
				"ANY /* - All traffic proxied to tunnel",
			},
		})
	})

	// WebSocket endpoint for tunnel client (ONLY management route)
	r.GET("/_tunnel/ws/connect", tunnelHandlers.HandleWebSocket)

	// Catch-all for ALL traffic (proxy everything else)
	r.NoRoute(tunnelHandlers.HandleTunnelTraffic)

	// Start server
	log.Printf("Starting Rail Tunnel server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
*/
