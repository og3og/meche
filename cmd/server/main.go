package main

import (
	"meche/pkg/config"
)

func main() {
	// Create and configure the server
	server := config.NewServer()

	// Start server
	server.Logger.Fatal(server.Start(":3000"))
}
