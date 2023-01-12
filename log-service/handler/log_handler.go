package handler

import (
	"log_service/entity"
	"log_service/repo"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type JsonResponse struct {
	Error   bool
	Message string
	Data    any `json:",omitempty"`
}

type LogHandler interface {
	InsertData(c echo.Context) error
}

type LogHandlerImpl struct {
	Repo repo.MongoRepo
}

func NewLogHandler(repo repo.MongoRepo) LogHandler {
	return &LogHandlerImpl{
		Repo: repo,
	}
}

func (handler *LogHandlerImpl) InsertData(c echo.Context) error {
	model := new(entity.PayloadModel)
	err := c.Bind(model)
	if err != nil {
		return c.JSON(http.StatusBadRequest, JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	err = handler.Repo.InsertData(&entity.LogEntry{
		Name:      model.Name,
		Data:      model.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, JsonResponse{
		Error:   false,
		Message: "success logged data",
	})

}
