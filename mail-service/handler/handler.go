package handler

import (
	"mail_service/entity"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
)

type JsonResponse struct {
	Error   bool
	Message string
	Data    any `json:",omitempty"`
}

type MailHandler interface {
	SendMail(c echo.Context) error
}

type MailHandlerImpl struct{}

func (handler *MailHandlerImpl) SendMail(c echo.Context) error {
	reqModel := new(entity.MailPayload)
	err := c.Bind(reqModel)
	if err != nil {
		return c.JSON(http.StatusBadRequest, JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	SMTP_HOST := os.Getenv("SMTP_HOST")
	SMTP_PORT, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	SENDER := "Mail Service NeoJ <noreply-ocra@neojarma.com>"
	SMTP_USER := os.Getenv("SMTP_USER")
	SMTP_PASS := os.Getenv("SMTP_PASS")

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", SENDER)
	mailer.SetHeader("To", reqModel.To)
	mailer.SetHeader("Subject", "Test mail")
	mailer.SetBody("text/html", reqModel.Message)

	dialer := gomail.NewDialer(
		SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASS,
	)

	if err := dialer.DialAndSend(mailer); err != nil {
		return c.JSON(http.StatusInternalServerError, JsonResponse{
			Error:   true,
			Message: "failed to send mail : " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, JsonResponse{
		Error:   false,
		Message: "mail has been sent..",
	})

}
