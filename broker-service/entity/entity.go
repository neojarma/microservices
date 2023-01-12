package entity

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type JsonResponse struct {
	Error   bool
	Message string
	Data    any `json:",omitempty"`
}

type ActionPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MailPayload struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

type Config struct {
	RabbitConn *amqp.Connection
	RabbitChan *amqp.Channel
}
