package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type TCPChecker struct {
	IP      net.IP
	Port    uint16
	Timeout time.Duration
}

type Result struct {
	Success bool
	Message string
}

func NewTCPChecker(ip net.IP, port uint16, timeout time.Duration) *TCPChecker {
	return &TCPChecker{
		ip,
		port,
		timeout,
	}
}

func (tcpChecker *TCPChecker) Check(timeout time.Duration) *Result {
	connStartTime := time.Now()
	conn, err := net.DialTimeout("tcp", tcpChecker.addr(), timeout)
	if err != nil {
		return &Result{
			false,
			fmt.Sprintf("Failed to connect: %v: ", err),
		}
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %s", err)
		}
		log.Printf("Connection closed: %s. Connection duration: %v ms", conn.RemoteAddr(), time.Since(connStartTime))
	}()

	buf := make([]byte, 1024)
	for {
		size, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf("Error reading from connection: %s", err)
		} else {
			log.Printf("Read %d bytes: %s", size, string(buf[:size]))
		}
		break
	}

	return &Result{
		true,
		fmt.Sprintf("Connected to %s:%d", tcpChecker.IP, tcpChecker.Port),
	}
}

func (tcpChecker *TCPChecker) CheckWithRetries(retries int, retryDelay time.Duration) *Result {
	var result *Result
	for i := 0; i < retries; i++ {
		start := time.Now()
		result = tcpChecker.Check(tcpChecker.Timeout)
		duration := time.Since(start)
		fmt.Printf("Health check attempt %d - Success: %v, Latency: %v, Message: %s\n", i+1, result.Success, duration, result.Message)

		if result.Success {
			return result
		}
		time.Sleep(retryDelay)
	}
	return result
}

func (tcpChecker *TCPChecker) addr() string {
	return fmt.Sprintf("%s:%d", tcpChecker.IP, tcpChecker.Port)
}
