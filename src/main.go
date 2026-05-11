package main

import (
	"fmt"
	"log"
	"net/http"
	
	"aether-db/server"
	"aether-db/storage"
)

const (
	DefaultPort = ":8080"
	Version     = "1.0.0"
)

func main() {
	fmt.Println("🗄️  AetherDB - Distributed Key-Value Store")
	fmt.Printf("   Version: %s\n\n", Version)
	
	// Initialize storage engine
	db := storage.NewStore()
	defer db.Close()
	
	// Create HTTP server
	srv := server.NewServer(db)
	
	// Routes
	http.HandleFunc("/", srv.HandleRoot)
	http.HandleFunc("/set", srv.HandleSet)
	http.HandleFunc("/get", srv.HandleGet)
	http.HandleFunc("/delete", srv.HandleDelete)
	http.HandleFunc("/keys", srv.HandleKeys)
	http.HandleFunc("/health", srv.HandleHealth)
	
	// Start server
	fmt.Printf("Server listening on %s\n", DefaultPort)
	log.Fatal(http.ListenAndServe(DefaultPort, nil))
}
