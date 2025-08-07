# 🚇 Rail Tunnel Server

> **Professional tunneling server that exposes local development servers to the internet via WebSocket connections**

Rail Tunnel Server is the backend component of the Rail Tunnel ecosystem - a production-ready tunneling solution that allows developers to securely expose their local development servers to the internet. Built with Go and designed for cloud deployment.

## 🌟 Features

- **WebSocket Tunneling** - Real-time bidirectional communication
- **HTTP Request Proxying** - Forward HTTP requests to connected clients
- **Multiple Client Support** - Handle multiple tunnel connections simultaneously
- **Health Monitoring** - Built-in health check and status endpoints
- **Cloud Ready** - Optimized for Railway, Heroku, and Docker deployment
- **Graceful Shutdown** - Proper cleanup and connection handling
- **Production Logging** - Structured logging with Gin framework

## 🚀 Quick Start

### Prerequisites
- Go 1.21+ installed
- Git for cloning

### Local Development

```bash
# Clone the repository
git clone https://github.com/isaui/rail-tunnel.git
cd rail-tunnel

# Install dependencies
go mod tidy

# Run the server
go run main.go

# Server starts on port 8080 (or PORT env variable)
```

### Using with Rail Tunnel CLI

1. **Install the CLI**:
   ```bash
   npm install -g rail-tunnel
   ```

2. **Start your local app** (e.g., on port 3000)

3. **Connect via tunnel**:
   ```bash
   npx rail-tunnel tunnel --port 3000 --remote http://localhost:8080
   ```

4. **Access your app** via the tunnel server URL

## 🔗 API Endpoints

### Core Endpoints
- `GET /` - Service information and status
- `GET /_tunnel/health` - Health check endpoint
- `GET /_tunnel/info` - Tunnel server information
- `GET /ws/connect` - WebSocket connection for tunnel clients
- `ALL /*` - Proxy requests to connected tunnel clients

### Health Check Response
```json
{
  "status": "healthy",
  "service": "rail-tunnel-server",
  "timestamp": "2025-01-07T21:30:00Z"
}
```

## 🐳 Docker Deployment

### Build Docker Image
```bash
# Build the image
docker build -t rail-tunnel-server .

# Run the container
docker run -p 8080:8080 rail-tunnel-server
```

### Docker Compose
```yaml
version: '3.8'
services:
  rail-tunnel:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - GIN_MODE=release
```

## ⚙️ Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | `8080` | No |
| `GIN_MODE` | Gin framework mode | `debug` | No |

**Example:**
```bash
export PORT=9000
export GIN_MODE=release
go run main.go
```

## 🚀 Railway Deployment

### One-Click Deploy
[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/rail-tunnel)

### Manual Deployment

1. **Connect Repository**:
   - Go to [Railway Dashboard](https://railway.app/dashboard)
   - Click "New Project" → "Deploy from GitHub repo"
   - Select your forked `rail-tunnel` repository

2. **Configure Environment**:
   ```bash
   GIN_MODE=release
   # PORT is automatically set by Railway
   ```

3. **Deploy**:
   - Railway will automatically detect Go and build the project
   - Your tunnel server will be available at: `https://your-app.railway.app`

### Using Your Deployed Server

```bash
# Connect your local app via CLI
npx rail-tunnel tunnel --port 3000 --remote https://your-app.railway.app
```

## 🛠️ Development

### Project Structure
```
rail-tunnel/
├── main.go                 # Server entry point
├── handlers/
│   └── tunnel_handlers.go  # WebSocket & HTTP handlers
├── managers/
│   └── tunnel_manager.go   # Connection management
├── models/
│   └── tunnel.go          # Data structures
├── Dockerfile             # Docker configuration
└── go.mod                 # Go dependencies
```

### Testing

```bash
# Run server locally
go run main.go

# Test health endpoint
curl http://localhost:8080/_tunnel/health

# Test WebSocket connection (requires wscat)
wscat -c ws://localhost:8080/ws/connect
```

### Building

```bash
# Build binary
go build -o rail-tunnel main.go

# Run binary
./rail-tunnel
```

## 🔗 Related Projects

- **[Rail Tunnel CLI](https://github.com/isaui/rail-tunnel-cli)** - The companion CLI client
- **[NPM Package](https://www.npmjs.com/package/rail-tunnel)** - Install CLI via npm

## 🐛 Troubleshooting

### Common Issues

**"WebSocket connection failed"**
- Ensure server is running and accessible
- Check firewall settings
- Verify WebSocket endpoint: `ws://your-server/ws/connect`

**"Port already in use"**
- Change the PORT environment variable
- Kill existing process: `lsof -ti:8080 | xargs kill -9`

**"Build failed on Railway"**
- Ensure Go version compatibility (1.21+)
- Check go.mod and go.sum files are present
- Review Railway build logs for specific errors

## 📝 License

MIT - see [LICENSE](LICENSE) file for details

## 🤝 Contributing

Contributions are welcome! Please feel free to:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

**Built with ❤️ for the developer community**
