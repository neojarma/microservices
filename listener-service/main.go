package main

import (
	"bytes"
	"encoding/json"
	"listener_service/connection"
	"log"
	"net/http"
)

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func main() {
	rabbitConn, err := connection.GetConnection()
	if err != nil {
		log.Panic(err)
	}
	defer rabbitConn.Close()

	// chan
	ch, err := rabbitConn.Channel()
	if err != nil {
		log.Println("error open channel", err)
	}
	defer ch.Close()

	// queue
	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Println("error declare queue", err)
	}

	// exchange
	err = ch.ExchangeDeclare(
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
	}

	err = ch.QueueBind(
		q.Name,
		"log.*",
		"logs_topic",
		false,
		nil,
	)
	if err != nil {
		log.Println("error bind a queue", err)
	}

	// consume
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("error declare consumer", err)
	}

	for data := range msgs {
		var payload Payload
		err := json.Unmarshal(data.Body, &payload)
		if err != nil {
			log.Println("error unmarshal json", err)
		}

		handleBodyRequest(payload)
	}

}

func handleBodyRequest(payload Payload) {
	switch payload.Name {
	case "log":
		err := sendLogRequest(payload)
		if err != nil {
			log.Println("error sending log request from listener service", err)
		}
	default:
		err := sendLogRequest(payload)
		if err != nil {
			log.Println("error sending log request from listener service", err)
		}
	}

}

func sendLogRequest(data Payload) error {
	logHost := "http://logger-service/log"
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", logHost, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
