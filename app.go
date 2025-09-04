package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

// AppInfo represents information about the application
type AppInfo struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	GoVersion   string    `json:"go_version"`
	Platform    string    `json:"platform"`
	StartTime   time.Time `json:"start_time"`
	Uptime      string    `json:"uptime"`
	Environment string    `json:"environment"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime"`
}

// MetricsResponse represents basic application metrics
type MetricsResponse struct {
	RequestCount int64     `json:"request_count"`
	MemoryUsage  uint64    `json:"memory_usage_bytes"`
	Goroutines   int       `json:"goroutines"`
	Timestamp    time.Time `json:"timestamp"`
}

var (
	appStartTime  = time.Now()
	requestCount  int64
	appVersion    = "1.0.0"
	appName       = "Multi-Stage Docker Demo"
	environment   = getEnv("ENVIRONMENT", "development")
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func incrementRequestCount() {
	requestCount++
}

func getUptime() string {
	return time.Since(appStartTime).String()
}

// Root handler - welcome page
func rootHandler(w http.ResponseWriter, r *http.Request) {
	incrementRequestCount()
	
	response := map[string]interface{}{
		"message": "Welcome to the Multi-Stage Docker Demo!",
		"info":    "This application demonstrates the benefits of Docker multi-stage builds",
		"endpoints": map[string]string{
			"/":         "This welcome page",
			"/health":   "Health check endpoint",
			"/info":     "Application information",
			"/metrics":  "Basic application metrics",
			"/greet":    "Greeting service (POST with JSON: {\"name\": \"yourname\"})",
			"/demo":     "Docker build demo information",
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Health check handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	incrementRequestCount()
	
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    getUptime(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// App info handler
func infoHandler(w http.ResponseWriter, r *http.Request) {
	incrementRequestCount()
	
	info := AppInfo{
		Name:        appName,
		Version:     appVersion,
		GoVersion:   runtime.Version(),
		Platform:    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		StartTime:   appStartTime,
		Uptime:      getUptime(),
		Environment: environment,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// Metrics handler
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	incrementRequestCount()
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	metrics := MetricsResponse{
		RequestCount: requestCount,
		MemoryUsage:  m.Sys,
		Goroutines:   runtime.NumGoroutine(),
		Timestamp:    time.Now(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// Demo handler - explains multi-stage builds
func demoHandler(w http.ResponseWriter, r *http.Request) {
	incrementRequestCount()
	
	demo := map[string]interface{}{
		"title": "Docker Multi-Stage Build Benefits",
		"benefits": []string{
			"Significantly smaller image sizes",
			"Improved security (no build tools in production)",
			"Better layer caching and optimization",
			"Separation of build and runtime environments",
			"Reduced attack surface",
		},
		"comparison": map[string]interface{}{
			"single_stage": map[string]string{
				"description": "Traditional approach - includes all build tools",
				"typical_size": "~1GB+ (includes Go compiler, tools, source)",
				"security": "Lower (contains build tools and dependencies)",
			},
			"multi_stage": map[string]string{
				"description": "Optimized approach - only runtime binary",
				"typical_size": "~10-50MB (only binary and minimal OS)",
				"security": "Higher (minimal attack surface)",
			},
		},
		"implementation": map[string]string{
			"build_stage": "Uses full golang image for compilation",
			"runtime_stage": "Uses minimal base image (alpine/scratch)",
			"optimization": "Static linking, stripped binaries",
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(demo)
}

// Greet handler - enhanced greeting service
func greetHandler(w http.ResponseWriter, r *http.Request) {
	incrementRequestCount()
	
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed. Use POST with JSON payload: {\"name\": \"yourname\"}",
		})
		return
	}
	
	var req struct {
		Name string `json:"name"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid JSON payload",
		})
		return
	}
	
	if req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Name field is required",
		})
		return
	}
	
	response := map[string]interface{}{
		"message":     fmt.Sprintf("Hello, %s! Welcome to the Multi-Stage Docker Demo!", req.Name),
		"timestamp":   time.Now(),
		"environment": environment,
		"served_by":   fmt.Sprintf("%s v%s", appName, appVersion),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Logging middleware
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	}
}

// Simple router
func route(pattern string, handler http.HandlerFunc) {
	http.HandleFunc(pattern, loggingMiddleware(handler))
}

func main() {
	// Get port from environment or use default
	port := getEnv("PORT", "8080")
	
	// Routes
	route("/", rootHandler)
	route("/health", healthHandler)
	route("/info", infoHandler)
	route("/metrics", metricsHandler)
	route("/greet", greetHandler)
	route("/demo", demoHandler)
	
	// Start server
	log.Printf("Starting %s v%s on port %s", appName, appVersion, port)
	log.Printf("Environment: %s", environment)
	log.Printf("Go version: %s", runtime.Version())
	log.Printf("Platform: %s/%s", runtime.GOOS, runtime.GOARCH)
	
	addr := ":" + port
	
	// For Docker containers, listen on all interfaces
	if strings.Contains(environment, "production") || os.Getenv("DOCKER") == "true" {
		addr = "0.0.0.0:" + port
	}
	
	log.Printf("Server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
