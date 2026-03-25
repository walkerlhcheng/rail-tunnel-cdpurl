package main

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
				"WS /?port=3000 - WebSocket connection for tunnel client (Moved to root)",
				"ANY /* - All traffic proxied to tunnel",
			},
		})
	})

	// ==========================================
	// 門神 (The Bouncer) - 處理大門口 ( / ) 嘅分流
	// ==========================================
	r.GET("/", func(c *gin.Context) {
		// 只要係連去根目錄 (/) 嘅 WebSocket，就當係 Tunnel CLI Client
		if c.IsWebsocket() {
			// Hardcode 設定 port 為 3000，等 tunnelHandlers 認到
			c.Request.URL.RawQuery = "port=3000"
			tunnelHandlers.HandleWebSocket(c)
			return
		}
	// 其他普通 GET 請求 (唔係 WebSocket) -> 交畀 Proxy 轉發
		tunnelHandlers.HandleTunnelTraffic(c)
	})

	/*
	// ==========================================
	// 門神 (The Bouncer) - 處理大門口 ( / ) 嘅分流
	// ==========================================
	r.GET("/", func(c *gin.Context) {
		// 檢查係咪 WebSocket Upgrade，同埋有冇帶 ?port= 參數
		if c.IsWebsocket() && c.Query("port") != "" {
			// 認得係自己友 (Tunnel CLI Client) -> 建立骨幹隧道
			tunnelHandlers.HandleWebSocket(c)
			return
		}

		// 其他普通 GET 請求 (冇 port 參數或者唔係 WebSocket) -> 交畀 Proxy 轉發
		tunnelHandlers.HandleTunnelTraffic(c)
	})
	*/

	// ==========================================
	// Proxy Catch-all 規則
	// ==========================================
	// 處理所有其他路徑 (例如 /devtools/browser/XXXX)
	r.NoRoute(tunnelHandlers.HandleTunnelTraffic)
	// 處理大門口 ( / ) 嘅其他 Method (例如 POST, PUT 等等)
	r.NoMethod(tunnelHandlers.HandleTunnelTraffic)

	// Start server
	log.Printf("Starting Rail Tunnel server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
