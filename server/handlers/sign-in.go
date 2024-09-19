package handlers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quyld17/chat-app/entities/users"
	jwtService "github.com/quyld17/chat-app/services/jwt"
)

func SignIn(c echo.Context, db *sql.DB) error {
	var account users.Users
	if err := c.Bind(&account); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err := users.Authenticate(account, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	userId, err := users.GetIdByUsername(c, db, account.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	token, err := jwtService.Generate(account.Username, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, echo.Map{"token": token})
}
