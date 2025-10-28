package routers

import (
	"database/sql"

	"auth/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterAPIHandlers(router *echo.Echo, dbMySQL *sql.DB) {
	router.POST("/sign-up", func(c echo.Context) error {
		return handlers.SignUp(c, dbMySQL)
	})
	router.POST("/sign-in", func(c echo.Context) error {
		return handlers.SignIn(c, dbMySQL)
	})
}
