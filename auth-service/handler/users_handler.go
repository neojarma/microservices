package handler

import (
	"auth_service/entity"
	"auth_service/repo"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type JsonResponse struct {
	Error   bool
	Message string
	Data    any `json:",omitempty"`
}

type UserHandler interface {
	Authenticate(c echo.Context) error
	IsPasswdValid(ashedPasswd, plainPasswd string) (bool, error)
}

type UserHandlerImpl struct {
	Repo repo.UserRepo
}

func NewUserHandler(repo repo.UserRepo) UserHandler {
	return &UserHandlerImpl{
		Repo: repo,
	}
}

func (handler *UserHandlerImpl) Authenticate(c echo.Context) error {
	payload := new(entity.AuthPayload)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusUnauthorized, JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	user, err := handler.Repo.GetByEmail(payload.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	// check pass
	valid, err := handler.IsPasswdValid(user.Password, payload.Password)
	if !valid || err != nil {
		return c.JSON(http.StatusUnauthorized, JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	err = callLogService(&LogPayload{
		Name: "auth log",
		Data: fmt.Sprintf("%s is logged in", user.Email),
	}, c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, JsonResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, JsonResponse{
		Error:   false,
		Message: "authenticated",
		Data:    user,
	})

}

func callLogService(auth *LogPayload, ctx echo.Context) error {
	logHost := "http://logger-service/log"
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(auth)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, logHost, &buf)
	if err != nil {
		return err
	}

	resp, err := doHttpRequest(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = decodeJsonResponse(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

func (handler *UserHandlerImpl) IsPasswdValid(hashedPasswd, plainPasswd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPasswd), []byte(plainPasswd))

	if err != nil {
		return false, err
	}

	return true, nil
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

func decodeJsonResponse(bodyResponse io.ReadCloser) (*JsonResponse, error) {
	var jsonResponse JsonResponse
	err := json.NewDecoder(bodyResponse).Decode(&jsonResponse)
	if err != nil {
		return nil, err
	}

	return &jsonResponse, nil
}
