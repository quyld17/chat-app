package routers

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/quyld17/chat-app/handlers"
)

func RegisterAPIHandlers(router *echo.Echo, db *sql.DB) {
	//Authentication
	router.POST("/sign-up", func(c echo.Context) error {
		return handlers.SignUp(c, db)
	})
	router.POST("/sign-in", func(c echo.Context) error {
		return handlers.SignIn(c, db)
	})
}
