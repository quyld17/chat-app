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

	// WebSocket routes
	// router.GET("/ws/online", func(c echo.Context) error {
	// 	return handlers.OnlineStatus(c, db)
	// }
	router.GET("/ws/:conversation_id", func(c echo.Context) error {
		return handlers.ChatWebSocket(c, db)
	})
}
