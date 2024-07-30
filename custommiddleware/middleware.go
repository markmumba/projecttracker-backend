package custommiddleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markmumba/project-tracker/auth"
)


func Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("access_token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "unauthorized"})
		}

		claims, err := auth.ValidateAccessToken(cookie.Value)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "unauthorized"})
		}

		c.Set("userId", claims.UserId)

		return next(c)
	}
}