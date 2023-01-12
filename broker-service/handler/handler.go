package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"broker_service/entity"

	"github.com/labstack/echo/v4"
	"github.com/rabbitmq/amqp091-go"
)

var RabbitConfig = new(entity.Config)

const (
	AuthAction = "auth"
	LogAction  = "log"
	MailAction = "mail"
)

func Broker(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	return c.JSON(http.StatusOK, entity.JsonResponse{
		Error:   false,
		Message: "hit broker service",
	})
}

func HandleSubmission(c echo.Context) error {
	req := new(entity.ActionPayload)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, entity.JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	switch req.Action {
	case AuthAction:
		return callAuthService(&req.Auth, c)
	case LogAction:
		return callLogServiceRabbit(&req.Log, c)
	case MailAction:
		return callMailService(&req.Mail, c)
	default:
		return c.JSON(http.StatusBadRequest, entity.JsonResponse{
			Error:   true,
			Message: "unknown action",
		})
	}
}

func callMailService(auth *entity.MailPayload, ctx echo.Context) error {

	logHost := "http://mailer-service/mail"
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(auth)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, entity.JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	req, err := http.NewRequest(http.MethodPost, logHost, &buf)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, entity.JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	resp, err := doHttpRequest(req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, entity.JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	defer resp.Body.Close()

	jsonResponse, err := decodeJsonResponse(resp.Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, entity.JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	if resp.StatusCode != http.StatusOK {
		return ctx.JSON(http.StatusInternalServerError, entity.JsonResponse{
			Error:   jsonResponse.Error,
			Message: jsonResponse.Message,
		})
	}

	return ctx.JSON(resp.StatusCode, entity.JsonResponse{
		Error:   jsonResponse.Error,
		Message: jsonResponse.Message,
	})
}

func callLogServiceRabbit(data *entity.LogPayload, ctx echo.Context) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, entity.JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	err = RabbitConfig.RabbitChan.PublishWithContext(
		timeoutCtx,
		"logs_topic",
		"log.INFO",
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(jsonData),
		},
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, entity.JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, entity.JsonResponse{
		Error:   false,
		Message: "success log via rabbitmq",
	})
}

func callLogService(data *entity.LogPayload, ctx echo.Context) error {

	logHost := "http://logger-service/log"
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(data)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, entity.JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	req, err := http.NewRequest(http.MethodPost, logHost, &buf)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, entity.JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	resp, err := doHttpRequest(req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, entity.JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	defer resp.Body.Close()

	jsonResponse, err := decodeJsonResponse(resp.Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, entity.JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	if resp.StatusCode != http.StatusOK {
		return ctx.JSON(http.StatusInternalServerError, entity.JsonResponse{
			Error:   jsonResponse.Error,
			Message: jsonResponse.Message,
		})
	}

	return ctx.JSON(resp.StatusCode, entity.JsonResponse{
		Error:   jsonResponse.Error,
		Message: jsonResponse.Message,
	})
}

func callAuthService(auth *entity.AuthPayload, ctx echo.Context) error {
	authHost := "http://auth-service/auth"
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(auth)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, entity.JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	req, err := http.NewRequest(http.MethodPost, authHost, &buf)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, entity.JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	resp, err := doHttpRequest(req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, entity.JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	defer resp.Body.Close()

	jsonResponse, err := decodeJsonResponse(resp.Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, entity.JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	if resp.StatusCode != http.StatusOK {
		return ctx.JSON(http.StatusInternalServerError, entity.JsonResponse{
			Error:   jsonResponse.Error,
			Message: jsonResponse.Message,
		})
	}

	return ctx.JSON(resp.StatusCode, entity.JsonResponse{
		Error:   jsonResponse.Error,
		Message: jsonResponse.Message,
		Data:    jsonResponse.Data,
	})
}

func doHttpRequest(req *http.Request) (*http.Response, error) {
	// prevent http 415 Unsupported Media Type
	req.Header.Add("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func decodeJsonResponse(bodyResponse io.ReadCloser) (*entity.JsonResponse, error) {
	var jsonResponse entity.JsonResponse
	err := json.NewDecoder(bodyResponse).Decode(&jsonResponse)
	if err != nil {
		return nil, err
	}

	return &jsonResponse, nil
}
