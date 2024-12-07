package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// Parse the WebSocket URL from command-line arguments
	var midServerAddr string
	flag.StringVar(&midServerAddr, "mid", "ws://localhost:8100/connect", "WebSocket address of the Mid-Server")
	flag.Parse()

	// Validate the WebSocket URL
	u, err := url.Parse(midServerAddr)
	if err != nil {
		log.Fatalf("Invalid WebSocket address: %v", err)
	}

	// Establish WebSocket connection to Mid-Server
	log.Printf("Connecting to Mid-Server at %s...\n", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("Failed to connect to Mid-Server: %v", err)
	}
	defer conn.Close()
	log.Println("Connected to Mid-Server via WebSocket.")

	// Handle interrupt signal to cleanly close the connection
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Start a ticker to send heartbeat messages every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Main loop: send heartbeat and read messages
	for {
		select {
		case <-ticker.C:
			// Send a heartbeat message to Mid-Server
			heartbeat := "Heartbeat from Rad-Server"
			if err := conn.WriteMessage(websocket.TextMessage, []byte(heartbeat)); err != nil {
				log.Printf("Failed to send heartbeat: %v", err)
				return
			}
			log.Println("Sent heartbeat to Mid-Server.")

		case <-interrupt:
			// Gracefully close the connection on interrupt signal
			log.Println("Interrupt signal received. Closing WebSocket connection.")
			if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Rad-Server shutting down")); err != nil {
				log.Printf("Error closing connection: %v", err)
			}
			return

		default:
			// Read messages from the Mid-Server
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				return
			}
			log.Printf("Message from Mid-Server: %s", string(message))
		}
	}
}
