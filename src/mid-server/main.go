package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/armon/go-socks5"
	"github.com/quic-go/quic-go"
)

var radConnection quic.Connection // Global QUIC connection to Rad-Server

func main() {
	// Start SOCKS5 proxy
	go startSocks5Proxy()

	// Accept Rad-Server connection
	acceptRadConnection()

	select {} // Keep the main thread alive
}

// acceptRadConnection accepts a QUIC connection from Rad-Server
func acceptRadConnection() {
	listener, err := quic.ListenAddr("0.0.0.0:8100", generateTLSConfig(), nil)
	if err != nil {
		log.Fatalf("Failed to start QUIC listener: %v", err)
	}
	log.Println("Listening for Rad-Server connections on 0.0.0.0:8100...")

	// Accept the Rad-Server connection
	connection, err := listener.Accept(context.Background())
	if err != nil {
		log.Fatalf("Failed to accept Rad-Server connection: %v", err)
	}
	radConnection = connection
	log.Println("Rad-Server connected.")
}

// startSocks5Proxy starts the SOCKS5 proxy server
func startSocks5Proxy() {
	conf := &socks5.Config{
		Dial: customDialToRadServer, // Custom dialer to forward traffic via Rad-Server
	}
	server, err := socks5.New(conf)
	if err != nil {
		log.Fatalf("Failed to create SOCKS5 server: %v", err)
	}

	address := "0.0.0.0:1080"
	log.Printf("SOCKS5 Proxy Server listening on %s\n", address)
	if err := server.ListenAndServe("tcp", address); err != nil {
		log.Fatalf("Failed to start SOCKS5 server: %v", err)
	}
}

// customDialToRadServer forwards SOCKS5 traffic to Rad-Server using QUIC streams
func customDialToRadServer(ctx context.Context, network, addr string) (net.Conn, error) {
	if radConnection == nil {
		return nil, fmt.Errorf("No connection to Rad-Server established")
	}

	// Open a new QUIC stream to Rad-Server
	stream, err := radConnection.OpenStreamSync(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to open QUIC stream to Rad-Server: %v", err)
	}

	// Send the CONNECT request (target address) to Rad-Server
	_, err = stream.Write([]byte(addr))
	if err != nil {
		return nil, fmt.Errorf("Failed to send CONNECT request to Rad-Server: %v", err)
	}

	log.Printf("Forwarding request to Rad-Server: %s", addr)

	// Return the QUIC stream as a net.Conn for SOCKS5 compatibility
	return &quicStreamWrapper{stream}, nil
}

// quicStreamWrapper adapts a QUIC stream to net.Conn
type quicStreamWrapper struct {
	quic.Stream
}

func (q *quicStreamWrapper) LocalAddr() net.Addr  { return nil }
func (q *quicStreamWrapper) RemoteAddr() net.Addr { return nil }
func (q *quicStreamWrapper) SetDeadline(t time.Time) error {
	return nil
}
func (q *quicStreamWrapper) SetReadDeadline(t time.Time) error {
	return nil
}
func (q *quicStreamWrapper) SetWriteDeadline(t time.Time) error {
	return nil
}

// generateTLSConfig generates the TLS configuration for QUIC
func generateTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true, // Skip verification for simplicity
		NextProtos:         []string{"quic-go"},
	}
}
