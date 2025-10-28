package routers

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"auth/handlers"
)

func RegisterAPIHandlers(router *echo.Echo, dbMySQL *sql.DB) {
	router.POST("/sign-up", func(c echo.Context) error {
		return handlers.SignUp(c, dbMySQL)
	})
	router.POST("/sign-in", func(c echo.Context) error {
		return handlers.SignIn(c, dbMySQL)
	})
}
