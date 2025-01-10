package controllers

import (
	"jan540/save-state/db"
	"jan540/save-state/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	db        *db.SaveDB
	jwtSecret string
}

func NewAuthController(db *db.SaveDB, jwtSecret string) *AuthController {
	return &AuthController{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

type LoginReq struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

type LoginRes struct {
	Token string `json:"token"`
}

func (ac *AuthController) Login(c echo.Context) error {
	var req LoginReq

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	user, err := ac.db.GetUserPassword(req.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid username :(")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid password :(")
	}

	claims := &jwt.RegisteredClaims{
		Issuer:    "save-state",
		Subject:   user.UserId,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 365)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(ac.jwtSecret))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create session"+err.Error())
	}

	return c.JSON(http.StatusOK, LoginRes{Token: signedToken})
}

type RegisterReq struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

func (ac *AuthController) Register(c echo.Context) error {
	var req RegisterReq

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	exists, err := ac.db.UserExists(req.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check if user exists")
	}

	if exists {
		return echo.NewHTTPError(http.StatusBadRequest, "User already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
	}

	user := &models.User{
		UserId:   uuid.NewString(),
		Username: req.Username,
		Password: string(hashedPassword),
	}

	err = ac.db.CreateUser(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	return c.JSON(http.StatusOK, "User created")
}
