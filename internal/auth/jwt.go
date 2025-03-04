package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(userId uint64, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userId,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

}
