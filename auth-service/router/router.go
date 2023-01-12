package router

import (
	"auth_service/handler"
	"auth_service/repo"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func Router(db *gorm.DB) *echo.Echo {
	app := echo.New()

	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXCSRFToken, echo.HeaderAuthorization},
		AllowMethods:     []string{echo.GET, echo.PUT, echo.DELETE, echo.POST, echo.OPTIONS},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	userRepo := repo.NewUserRepo(db)
	userHandler := handler.NewUserHandler(userRepo)
	app.POST("/auth", userHandler.Authenticate)

	return app
}
