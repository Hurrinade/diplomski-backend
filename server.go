package main

import (
	"github.com/Hurrinade/diplomski-backend/router"
)

func main() {
	// Create service which implements grpc interface and has auth service for firebase authentication
	r := router.NewRouter()

	r.Run(":8080")
}