package routers

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/quyld17/chat-app/handlers"
	"github.com/quyld17/chat-app/middlewares"
	"github.com/redis/go-redis/v9"
)

func RegisterAPIHandlers(router *echo.Echo, dbMySQL *sql.DB, dbRedis *redis.Client) {
	router.POST("/sign-up", func(c echo.Context) error {
		return handlers.SignUp(c, dbMySQL)
	})
	router.POST("/sign-in", func(c echo.Context) error {
		return handlers.SignIn(c, dbMySQL)
	})
	router.POST("/google-sign-in", func(c echo.Context) error {
		return handlers.GoogleSignIn(c, dbMySQL)
	})

	router.GET("/online-list", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.GetOnlineList(c, dbMySQL)
	}))

	router.GET("/ws/status", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.UpdateStatus(c, dbMySQL)
	}))

	router.GET("/chat-history", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.GetChatHistory(c, dbMySQL)
	}))

	router.GET("/ws/chat", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.Chat(c, dbMySQL, dbRedis)
	}))
}
