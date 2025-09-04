# Docker Multi-Stage Build Demo

This repository demonstrates the powerful benefits of Docker multi-stage builds using a Go web application. The demo showcases how multi-stage builds can dramatically reduce image sizes, improve security, and optimize deployment packages.

## 🚀 What This Demo Shows

This enhanced demo includes:

- **Comprehensive Web Application**: A feature-rich Go web server with multiple endpoints
- **Educational Content**: Built-in `/demo` endpoint explaining multi-stage benefits
- **Real-world Comparison**: Side-by-side comparison of single-stage vs multi-stage builds
- **Best Practices**: Optimized Dockerfiles demonstrating production-ready techniques

## 📊 Size Comparison

| Build Type | Image Size | Reduction Factor |
|------------|------------|------------------|
| Single-Stage | ~886MB | Baseline |
| Multi-Stage | ~4.81MB | **184x smaller** |

## 🏗️ Application Features

The Go application includes these endpoints:

- `GET /` - Welcome page with endpoint overview
- `GET /health` - Health check endpoint with uptime
- `GET /info` - Application information (version, platform, etc.)
- `GET /metrics` - Basic runtime metrics (memory, goroutines, requests)
- `POST /greet` - Greeting service (send JSON: `{"name": "yourname"}`)
- `GET /demo` - **Educational content about multi-stage builds**

## 🔧 How to Build and Run

### Prerequisites

- Docker installed
- Go 1.21+ (for local development)

### Building the Images

```bash
# Build single-stage image (traditional approach)
docker build -f Dockerfile.single -t golang-single .

# Build multi-stage image (optimized approach)
docker build -f Dockerfile.multi -t golang-multi .
```

### Running the Applications

```bash
# Run single-stage version
docker run -p 8080:8080 golang-single

# Run multi-stage version
docker run -p 8080:8080 golang-multi
```

### Testing the Application

```bash
# Health check
curl http://localhost:8080/health

# Get application info
curl http://localhost:8080/info

# Learn about multi-stage builds
curl http://localhost:8080/demo

# Greeting service
curl -X POST -H "Content-Type: application/json" \
     -d '{"name": "Docker User"}' \
     http://localhost:8080/greet
```

## 📋 Multi-Stage Build Benefits

### 1. **Dramatic Size Reduction**
- **Single-stage**: Contains Go compiler, build tools, source code, and dependencies (~886MB)
- **Multi-stage**: Contains only the compiled binary and minimal runtime (~4.81MB)
- **Result**: 184x smaller images, faster deployments, lower storage costs

### 2. **Enhanced Security**
- **Build tools removed**: No Go compiler, git, or development tools in production
- **Minimal attack surface**: Only essential runtime components included
- **No source code**: Source files not included in final image

### 3. **Better Performance**
- **Faster deployments**: Smaller images download and start faster
- **Efficient caching**: Better Docker layer caching with separated build stages
- **Optimized binaries**: Static linking and stripped debug information

### 4. **Production Readiness**
- **Environment separation**: Clear distinction between build and runtime environments
- **Resource optimization**: Lower memory footprint in production
- **Container orchestration friendly**: Faster scaling with smaller images

## 🏗️ Technical Implementation

### Build Stage (Dockerfile.multi)
```dockerfile
# Uses full golang:1.21 image with all build tools
FROM golang:1.21 AS builder
WORKDIR /app
COPY app.go .
# Static linking for scratch compatibility
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s' -o app app.go
```

### Runtime Stage (Dockerfile.multi)
```dockerfile
# Uses minimal scratch image
FROM scratch
# Only copy the compiled binary
COPY --from=builder /app/app /app
ENV ENVIRONMENT=production
ENTRYPOINT ["/app"]
```

## 🔍 Inspecting the Images

Compare the image details:

```bash
# Check image sizes
docker images | grep golang

# Inspect image layers
docker history golang-single
docker history golang-multi

# Security scan comparison
docker scout quickview golang-single
docker scout quickview golang-multi
```

## 🧪 Development and Testing

### Local Development
```bash
# Run locally
go run app.go

# Build locally
go build app.go
./app
```

### Testing Endpoints
The application provides comprehensive information about itself and Docker multi-stage builds:

```bash
# Get complete demo information
curl http://localhost:8080/demo | jq .

# Check metrics
curl http://localhost:8080/metrics | jq .
```

## 🎯 Key Takeaways

1. **Multi-stage builds are essential** for production Go applications
2. **Size matters** - smaller images mean faster deployments and lower costs
3. **Security by design** - remove unnecessary tools from production images
4. **Performance optimization** - static linking and stripped binaries improve runtime
5. **Best practices** - separate build and runtime concerns for better maintainability

## 📚 Learn More

- [Docker Multi-stage builds documentation](https://docs.docker.com/build/building/multi-stage/)
- [Go Docker best practices](https://docs.docker.com/language/golang/)
- [Container security best practices](https://docs.docker.com/develop/dev-best-practices/)

## 🤝 Contributing

Feel free to submit issues and enhancement requests to improve this demo!

---

**This demo shows how a simple change in Dockerfile structure can result in 184x smaller images while improving security and performance.**