package main

import (
	"log"
	"net"
	"time"
)

type TCPListener struct {
	ip      string
	port    string
	message string
}

func NewTCPListener(ip string, port string, message string) *TCPListener {
	return &TCPListener{
		ip,
		port,
		message,
	}
}

func (listener *TCPListener) Listen() {
	address := net.JoinHostPort(listener.ip, listener.port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error listening on address %s: %s", address, err)
	}

	defer func() {
		if err := listen.Close(); err != nil {
			log.Fatalf("Error closing listener on address %s: %s", listen.Addr().String(), err)
		}
	}()

	log.Printf("Listening on address %s", listen.Addr().String())

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s", err)
			continue
		}
		go listener.handleConnection(conn, listener.message)
	}
}

func (listener *TCPListener) handleConnection(conn net.Conn, message string) {
	connStartTime := time.Now()
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %s", err)
		}
		log.Printf("Connection closed: %s. Connection duration: %v ms", conn.RemoteAddr(), time.Since(connStartTime))
	}()

	log.Printf("New connection from %s", conn.RemoteAddr())
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Printf("Error writing to connection: %s", err)
	}
}
