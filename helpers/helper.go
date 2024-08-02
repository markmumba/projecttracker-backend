package helpers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)


func ConvertUserID(c echo.Context, key string) (uuid.UUID, error) {
    userIDInterface := c.Get(key)
    switch id := userIDInterface.(type) {
    case string:
        // Try to parse the string as a UUID
        userID, err := uuid.Parse(id)
        if err != nil {
            return uuid.Nil, fmt.Errorf("invalid UUID format: %v", err)
        }
        return userID, nil
    default:
        return uuid.Nil, fmt.Errorf("invalid user ID type, expected string")
    }
}