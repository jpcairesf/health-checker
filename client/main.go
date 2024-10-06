package main

import (
	"net"
	"time"
)

func main() {
	checker := NewTCPChecker(net.ParseIP("127.0.0.1"), 8080, 10*time.Second)
	result := checker.CheckWithRetries(5, 2*time.Second)
	println("Result ", result.Message)
}
