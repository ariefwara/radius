package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/armon/go-socks5"
	"github.com/gorilla/websocket"
)

var (
	radServerConn *websocket.Conn // Global WebSocket connection
	wsMutex       sync.Mutex      // Mutex to prevent concurrent writes
)

func main() {
	// Start SOCKS5 Proxy and WebSocket Server
	go startSocks5Proxy()     // SOCKS5 server
	go startWebSocketServer() // WebSocket server

	// Keep the main thread alive
	select {}
}

// Start the SOCKS5 Proxy Server
func startSocks5Proxy() {
	conf := &socks5.Config{
		Rules: &loggingRule{}, // Custom rule for logging
	}
	server, err := socks5.New(conf)
	if err != nil {
		log.Fatalf("Failed to create SOCKS5 server: %v", err)
	}

	// Listen for SOCKS5 connections
	address := "0.0.0.0:1080"
	log.Printf("SOCKS Proxy Server listening on %s\n", address)
	if err := server.ListenAndServe("tcp", address); err != nil {
		log.Fatalf("Failed to start SOCKS5 server: %v", err)
	}
}

// Custom logging rule for SOCKS5 requests
type loggingRule struct{}

func (r *loggingRule) Allow(ctx context.Context, req *socks5.Request) (context.Context, bool) {
	log.Printf("New SOCKS request:\n  Command: %v\n  Destination: %s:%d\n",
		req.Command, req.DestAddr.IP, req.DestAddr.Port)

	// Forward data to the Rad-Server via WebSocket if connected
	if radServerConn != nil {
		message := fmt.Sprintf("SOCKS request to %s:%d", req.DestAddr.IP, req.DestAddr.Port)

		// Lock before writing to WebSocket
		wsMutex.Lock()
		defer wsMutex.Unlock()

		err := radServerConn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("Failed to send message to Rad-Server: %v", err)
		} else {
			log.Println("Forwarded request to Rad-Server via WebSocket")
		}
	} else {
		log.Println("No Rad-Server connected via WebSocket.")
	}

	return ctx, true
}

// Start WebSocket Server to Accept Rad-Server Connections
func startWebSocketServer() {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins for simplicity
	}

	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		var err error
		radServerConn, err = upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to upgrade to WebSocket: %v", err)
			return
		}
		defer radServerConn.Close()

		log.Println("Rad-Server connected via WebSocket.")
		for {
			// Read messages from Rad-Server
			_, message, err := radServerConn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket disconnected: %v", err)
				radServerConn = nil
				break
			}
			log.Printf("Message from Rad-Server: %s", string(message))
		}
	})

	address := "0.0.0.0:8100"
	log.Printf("WebSocket server listening on %s/connect\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("WebSocket server error: %v", err)
	}
}
