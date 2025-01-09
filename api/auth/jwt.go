package auth

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// TODO: dead code for now -> I only need UserId which is in sub
type JwtClaims struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GetUserIdFromContext(c echo.Context) (string, error) {
	user := c.Get("user").(*jwt.Token)
	userId, err := user.Claims.GetSubject()
	if err != nil {
		err = echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user id 🤷 HOW THE FUCK DID THIS HAPPEN?")
	}

	return userId, err
}
