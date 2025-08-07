# Rail Tunnel

Ngrok-like tunnel service that can be deployed on Railway.

## Features

- Health check endpoint
- Railway deployment ready
- Basic Gin server setup

## Quick Start

1. Install dependencies:
```bash
go mod tidy
```

2. Run the server:
```bash
go run main.go
```

3. Check health:
```bash
curl http://localhost:8080/health
```

## Endpoints

- `GET /` - Service information
- `GET /health` - Health check

## Environment Variables

- `PORT` - Server port (default: 8080)
- `GIN_MODE` - Gin mode (development/release)

## Railway Deployment

This service is ready to deploy on Railway. The `PORT` environment variable will be automatically set by Railway.
