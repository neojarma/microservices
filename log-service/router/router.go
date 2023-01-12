package router

import (
	"log_service/handler"
	"log_service/repo"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func Router(mongoClient *mongo.Client) *echo.Echo {
	app := echo.New()

	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXCSRFToken, echo.HeaderAuthorization},
		AllowMethods:     []string{echo.GET, echo.PUT, echo.DELETE, echo.POST, echo.OPTIONS},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	coll := mongoClient.Database("logs").Collection("logs")
	logRepo := repo.NewMongoRepo(coll)
	logHandler := handler.NewLogHandler(logRepo)
	app.POST("/log", logHandler.InsertData)

	return app
}
