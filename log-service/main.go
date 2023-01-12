package main

import (
	"context"
	"log"
	"log_service/connection"
	"log_service/router"
	"time"
)

func main() {
	client, err := connection.GetConnection()
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

	router := router.Router(client)
	router.Start(":80")
}
