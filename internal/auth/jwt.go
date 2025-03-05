package auth

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/macreai/chess-game-app-be/internal/entity"
	"github.com/macreai/chess-game-app-be/internal/model"
	"github.com/spf13/viper"
)

type MyJWT struct {
	User  *entity.User
	Viper *viper.Viper
}

func NewMyJWT(user *entity.User, viper *viper.Viper) *MyJWT {
	return &MyJWT{
		User:  user,
		Viper: viper,
	}
}

func (myJwt *MyJWT) GenerateJWT(user *entity.User, viper *viper.Viper) (string, error) {

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(viper.GetString("JWT_SECRET")))
}

func (myJwt *MyJWT) JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(&model.WebResponse[any]{
				Errors: fiber.NewError(fiber.StatusUnauthorized, "Missing or Malformed JWT"),
				Data:   nil,
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(&model.WebResponse[any]{
				Errors: fiber.NewError(fiber.StatusUnauthorized, "Invalid Auth Header"),
				Data:   nil,
			})
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected Signing Method")
			}
			return []byte(viper.GetString("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(&model.WebResponse[any]{
				Errors: fiber.NewError(fiber.StatusUnauthorized, "Invalid or Expired JWT"),
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("user_id", claims["user_id"])
			c.Locals("username", claims["username"])
		}

		return c.Next()
	}
}
