package main

import (
	"os"
	"github.com/wahyunurdian26/service-audit/transport"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	srv := transport.NewAMQPServer()
	
	go func() {
		// Use default port 8083 if unset
		port := os.Getenv("HTTP_PORT")
		if port == "" {
			port = "8083"
		}
		transport.RegisterHTTPServer(srv.Endpoints(), port)
	}()

	srv.Run()
}
