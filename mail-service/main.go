package main

import "mail_service/router"

func main() {
	router := router.Router()

	router.Start(":80")
}
