package routers

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/quyld17/chat-app/handlers"
	"github.com/quyld17/chat-app/middlewares"
)

func RegisterAPIHandlers(router *echo.Echo, db *sql.DB) {
	router.POST("/sign-up", func(c echo.Context) error {
		return handlers.SignUp(c, db)
	})
	router.POST("/sign-in", func(c echo.Context) error {
		return handlers.SignIn(c, db)
	})

	router.GET("/ws/status", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.UpdateStatus(c, db)
	}))

	router.GET("/online-list", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.GetOnlineList(c, db)
	}))

	router.GET("/ws/chat", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.Chat(c, db)
	}))
}
