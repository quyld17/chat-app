package handlers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/joho/godotenv"
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
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	userId, err := users.GetIdByUsername(c, db, account.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	token, err := jwtService.Generate(account.Username, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func GoogleSignIn(c echo.Context, db *sql.DB) error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to retrieve credentials. Please try again")
	}
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientId == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get Google Client ID")
	}

	var googleToken string
	if err := c.Bind(&googleToken); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	ctx := context.Background()
	payload, err := idtoken.Validate(ctx, googleToken, googleClientId)
	if err != nil {
		log.Printf("Token verification failed: %v", err)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	email := payload.Claims["email"].(string)
	err = users.CheckOrCreateGoogleAccount(c, db, email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	userId, err := users.GetIdByUsername(c, db, email)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	token, err := jwtService.Generate(email, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, echo.Map{"token": token})
}
