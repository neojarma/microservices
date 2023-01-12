package router

import (
	"broker_service/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Router() *echo.Echo {
	app := echo.New()

	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXCSRFToken, echo.HeaderAuthorization},
		AllowMethods:     []string{echo.GET, echo.PUT, echo.DELETE, echo.POST, echo.OPTIONS},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	app.POST("/", handler.Broker)
	app.POST("/handle", handler.HandleSubmission)

	return app
}
