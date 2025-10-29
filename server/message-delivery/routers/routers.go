package routers

import (
	"message-delivery/handlers"
	"message-delivery/services/jwt"

	"github.com/labstack/echo/v4"
)

func RegisterAPIHandlers(router *echo.Echo) {
	router.GET("/ws/setup", jwt.Authorize(func(c echo.Context) error {
		return handlers.SetUpConnection(c)
	}))
}
