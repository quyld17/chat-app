package handlers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quyld17/chat-app/entities/users"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c echo.Context, dbMySQL *sql.DB) error {
	var newUser users.Users
	if err := c.Bind(&newUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create new account! Please try again")
	}
	newUser.Password = string(hashedPassword)

	if err := users.Create(newUser, dbMySQL); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Account already exists! Please try again.")
	}

	return c.JSON(http.StatusOK, "Account created successfully!")
}
