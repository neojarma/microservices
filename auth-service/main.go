package main

import (
	"auth_service/connection"
	"auth_service/router"
	"log"
)

func main() {
	// initialize db
	db, err := connection.GetConnection()
	if err != nil {
		log.Fatalln(err.Error())
	}

	// start web server
	app := router.Router(db)
	app.Start(":80")
}
