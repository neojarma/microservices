package main

import (
	"broker_service/entity"
	"broker_service/handler"
	"broker_service/router"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	app := router.Router()

	rabbitConfig, err := setupRabbit()
	if err != nil {
		log.Println(err)
	}

	handler.RabbitConfig = rabbitConfig
	defer rabbitConfig.RabbitConn.Close()

	app.Start(":80")
}

func setupRabbit() (*entity.Config, error) {
	connection, err := connect()
	if err != nil {
		log.Println("error connecting to rabbitmq", err)
		return nil, err
	}

	rabbitChan, err := connection.Channel()
	if err != nil {
		log.Println("error open a channel from broker service", err)
		return nil, err
	}

	err = rabbitChan.ExchangeDeclare(
		"logs_topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("error declare an exchange", err)
		return nil, err
	}

	return &entity.Config{
		RabbitConn: connection,
		RabbitChan: rabbitChan,
	}, nil
}

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
