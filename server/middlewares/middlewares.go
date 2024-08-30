package middlewares

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	jwtService "github.com/quyld17/chat-app/services/jwt"
)

func JWTAuthorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := godotenv.Load(".env")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		tokenString := jwtService.GetToken(c)
		if tokenString == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})
		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		username := jwtService.GetClaims(token, "username")
		if username == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid claims")
		}
		c.Set("username", username)

		return next(c)
	}
}

// func Pagination(c echo.Context, itemsPerPage int) (int, error) {
// 	pageStr := c.QueryParam("page")
// 	page, err := strconv.Atoi(pageStr)
// 	if err != nil {
// 		return 0, fmt.Errorf("Invalid page number")
// 	}

// 	offset := (page - 1) * itemsPerPage

// 	return offset, nil
// }
