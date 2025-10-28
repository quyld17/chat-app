package main

import (
	"net/http"
	"time"

	"message-delivery/routers"
	"message-delivery/services/databases"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	dbMySQL := databases.NewMySQL()

	router := echo.New()
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		ExposeHeaders:    []string{echo.HeaderContentLength},
		AllowCredentials: true,
		MaxAge:           int(24 * time.Hour.Seconds()),
	}))

	routers.RegisterAPIHandlers(router, dbMySQL)

	router.Logger.Fatal(router.Start(":8080"))
}
