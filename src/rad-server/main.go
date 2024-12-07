package main

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"
	"strings"

	"github.com/quic-go/quic-go"
)

const midServerAddr = "localhost:8100" // Mid-Server address

func main() {
	// Connect to Mid-Server
	conn, err := quic.DialAddr(context.Background(), midServerAddr, generateTLSConfig(), nil)
	if err != nil {
		log.Fatalf("Failed to connect to Mid-Server: %v", err)
	}
	log.Println("Connected to Mid-Server.")

	// Continuously listen for forwarded requests from Mid-Server
	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			log.Printf("Error accepting stream: %v. Retrying...", err)
			continue
		}
		log.Println("New QUIC stream accepted.")
		go handleStream(stream) // Handle each request in its own goroutine
	}
}

// handleStream handles a forwarded CONNECT request from Mid-Server
func handleStream(stream quic.Stream) {
	defer stream.Close()

	// 1. Read the target address sent by Mid-Server
	buf := make([]byte, 1024)
	n, err := stream.Read(buf)
	if err != nil {
		log.Println("Error reading from Mid-Server:", err)
		return
	}
	target := strings.TrimSpace(string(buf[:n])) // Target address (e.g., example.com:443)
	log.Printf("Received CONNECT request for: %s", target)

	// 2. Connect to the actual target server
	conn, err := net.Dial("tcp", target)
	if err != nil {
		log.Printf("Failed to connect to target %s: %v", target, err)
		return
	}
	defer conn.Close()
	log.Printf("Connected to target: %s", target)

	// 3. Relay traffic bidirectionally
	relayTraffic(stream, conn)
}

// relayTraffic relays traffic between the QUIC stream and the target server
func relayTraffic(stream quic.Stream, conn net.Conn) {
	// Relay data from target → Mid-Server
	go func() {
		if _, err := io.Copy(stream, conn); err != nil {
			log.Println("Error relaying from target to Mid-Server:", err)
		}
		stream.Close()
	}()

	// Relay data from Mid-Server → target
	if _, err := io.Copy(conn, stream); err != nil {
		log.Println("Error relaying from Mid-Server to target:", err)
	}
	conn.Close()
	stream.Close()
}

// generateTLSConfig generates the TLS configuration for QUIC
func generateTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,       // Skip verification for simplicity
		ServerName:         "localhost", // Ensure this matches the expected SNI
		NextProtos:         []string{"quic-go"},
	}
}
