package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/markmumba/project-tracker/models"
	"golang.org/x/crypto/bcrypt"
)


var (
	accessTokenSecret  = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
	refreshTokenSecret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))
)

type JwtCustomClaims struct {
	UserId uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateAccessToken(user *models.User) (string, error) {
	stringId := user.ID.String() 

	claims := &JwtCustomClaims{
		user.ID,
		jwt.RegisteredClaims{
			Issuer:    stringId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)), // 15 minutes
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessTokenSecret)
}

func GenerateRefreshToken(user *models.User) (string, error) {
	stringId := user.ID.String() 

	claims := &JwtCustomClaims{
		user.ID,
		jwt.RegisteredClaims{
			Issuer:    stringId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 7 days
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshTokenSecret)
}

func ValidateAccessToken(tokenString string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return accessTokenSecret, nil
	})

	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func ValidateRefreshToken(tokenString string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return refreshTokenSecret, nil
	})

	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
