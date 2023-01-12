package connection

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func connect() (*amqp.Connection, error) {
	// specifying port :15672 or :5672 cause connection refused
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}

func GetConnection() (*amqp.Connection, error) {
	triesLimit := 5
	delay := time.Second * 5
	counter := 0

	for {
		conn, err := connect()
		if err != nil {
			log.Println("rabbitmq not yet ready...")
			counter++
		} else {
			log.Println("success to connect rabbitmq")
			return conn, nil
		}

		if counter >= triesLimit {
			log.Println("connection timed out: ", err)
			return nil, err
		}

		log.Print("still trying to connect...")
		time.Sleep(delay)
	}

}
