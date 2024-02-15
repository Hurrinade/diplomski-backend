package main

import (
	"context"
	"log"

	"github.com/Hurrinade/diplomski-backend/db"
	"github.com/Hurrinade/diplomski-backend/router"
)

func main() {
	// Create service which implements grpc interface and has auth service for firebase authentication
	cl := db.NewMongoClient()
	r := router.NewRouter(cl)

	r.Run(":8080")

	defer func() {
		log.Println("Closing db")
		if err := cl.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}