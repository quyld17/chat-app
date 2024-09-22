package handlers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quyld17/chat-app/entities/users"
)

func GetOnlineList(c echo.Context, db *sql.DB) error {
	userId := c.Get("user_id").(int)
	list, err := users.GetOnlineList(c, db, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve online users")
	}
	return c.JSON(http.StatusOK, list)
}
