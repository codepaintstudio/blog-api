package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"blog-api/pkg/config"
	"blog-api/pkg/database"
)

func main() {
	log.Println("Blog API Server starting...")
	
	// Initialize configuration
	_, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("Configuration loaded successfully")
	
	// Initialize database connections
	if err := database.InitMySQL(); err != nil {
		log.Fatalf("Failed to initialize MySQL: %v", err)
	}
	
	if err := database.InitRedis(); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	
	// Setup graceful shutdown
	setupGracefulShutdown()
	
	// TODO: Initialize routes and middleware
	// TODO: Start server
	
	log.Println("Blog API Server is ready")
	
	// Wait for shutdown signal
	waitForShutdown()
}

func setupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-c
		log.Println("Shutting down server...")
		
		// Close database connections
		if err := database.CloseMySQL(); err != nil {
			log.Printf("Error closing MySQL: %v", err)
		}
		
		if err := database.CloseRedis(); err != nil {
			log.Printf("Error closing Redis: %v", err)
		}
		
		log.Println("Server shutdown complete")
		os.Exit(0)
	}()
}

func waitForShutdown() {
	select {}
}