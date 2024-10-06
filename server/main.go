package main

func main() {
	NewTCPListener("127.0.0.1", "8080", "pong").Listen()
}
