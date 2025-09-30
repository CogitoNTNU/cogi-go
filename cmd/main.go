package main

import (
	"github.com/CogitoNTNU/cogi-go/internal/api"
)

// Entry point of the application
func main() {
	server, err := api.InitServer()

	defer func() {
		if r := recover(); r != nil {
			server.Logger.Info("Shutting down server...")
		}
	}()

	if (err != nil) {
		panic(err)
	}

	server.Serve()
}
