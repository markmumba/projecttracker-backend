package custommiddleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/markmumba/project-tracker/auth"
)

func Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "missing authorization header"})
		}

		// The Authorization header should be in the format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "invalid authorization header format"})
		}

		tokenString := parts[1]
		claims, err := auth.ValidateAccessToken(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "invalid or expired access token"})
		}

		// S.Println()
		c.Set("userId", claims.UserId)

		return next(c)
	}
}
