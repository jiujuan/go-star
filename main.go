package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	
	"go-star/pkg/graceful"
	"go-star/pkg/health"
)

func main() {
	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	// Create health manager
	healthManager := health.NewManager("go-star", "1.0.0")
	healthManager.SetTimeout(10 * time.Second)
	healthManager.SetCacheTTL(3 * time.Second)

	// Register basic health checks
	healthManager.RegisterFunc("application", func(ctx context.Context) health.CheckResult {
		return health.CheckResult{
			Name:      "application",
			Status:    health.StatusHealthy,
			Message:   "Application is running",
			Timestamp: time.Now(),
		}
	})

	// TODO: Add database and Redis health checks when those components are available
	// healthManager.Register("database", health.NewDatabaseChecker("database", db.Ping))
	// healthManager.Register("redis", health.NewRedisChecker("redis", redisClient.Ping))

	// Create graceful shutdown manager
	shutdownManager := graceful.NewManager()
	shutdownManager.SetTimeout(30 * time.Second)

	// Create Gin router
	router := gin.New()

	// Basic middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoints
	router.GET("/health", healthManager.HTTPHandler())
	router.GET("/health/ready", healthManager.ReadinessHandler())
	router.GET("/health/live", healthManager.LivenessHandler())

	// Legacy health endpoint for backward compatibility
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"service":   "go-star",
		})
	})

	// Default route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to go-star framework",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// Register HTTP server for graceful shutdown
	shutdownManager.Register("http-server", graceful.HTTPServerShutdown(srv))

	// Register cleanup tasks
	shutdownManager.Register("cleanup", graceful.CustomShutdown(func(ctx context.Context) error {
		log.Println("Performing cleanup tasks...")
		// Add any cleanup logic here
		return nil
	}))

	// Start server in a goroutine
	go func() {
		log.Printf("Starting go-star server on port %s", port)
		log.Printf("Health check endpoints:")
		log.Printf("  - Health: http://localhost:%s/health", port)
		log.Printf("  - Ready:  http://localhost:%s/health/ready", port)
		log.Printf("  - Live:   http://localhost:%s/health/live", port)
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for shutdown signal and perform graceful shutdown
	log.Println("Server started successfully. Press Ctrl+C to shutdown.")
	shutdownManager.Wait()
	log.Println("Server shutdown completed.")
}