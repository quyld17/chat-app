package handlers

import (
	"database/sql"
	"net/http"

	"auth/entities/users"
	jwtService "auth/services/jwt"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func SignIn(c echo.Context, dbMySQL *sql.DB) error {
	var account users.Users
	if err := c.Bind(&account); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err := users.Authenticate(account, dbMySQL)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	userId, err := users.GetIdByUsername(c, dbMySQL, account.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	token, err := jwtService.Generate(account.Username, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

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
