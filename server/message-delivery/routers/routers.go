package routers

import (
	"database/sql"
	"message-delivery/handlers"
	"message-delivery/services/jwt"

	"github.com/labstack/echo/v4"
)

func RegisterAPIHandlers(router *echo.Echo, dbMySQL *sql.DB) {
	router.GET("/ws/setup", jwt.Authorize(func(c echo.Context) error {
		return handlers.SetUpConnection(c)
	}))
}
