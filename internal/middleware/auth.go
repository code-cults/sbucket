package middleware

import (
	"os"
	"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == ""{
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error" : "No Authorization Header in Local Storage",
			})
		}
		
		parts := strings.Split(authHeader," ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error" : "Invalid header format",
			})
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface {}, error){
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized,"unexpected error")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error" : "invalid or expired token",
			})
		}
		claims := token.Claims.(jwt.MapClaims)
		userID := claims["sub"].(float64)

		c.Locals("userID",int(userID))
		return c.Next()
	}
}