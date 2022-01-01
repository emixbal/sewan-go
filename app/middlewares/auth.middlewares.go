package middlewares

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func IsAuthenticated(c *fiber.Ctx) error {
	raw_token := c.Request().Header.Peek("Authorization")
	tokenString := string(raw_token)

	if tokenString == "" {
		return c.Status(http.StatusUnauthorized).JSON(
			map[string]string{
				"message": "Unauthorized, need access token to access this API route!",
			},
		)
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		fmt.Println(err)
		return c.Status(http.StatusUnauthorized).JSON(
			map[string]string{
				"message": "Unauthorized, access token is invalid!",
			},
		)
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		c.Locals("user_id", claims["userId"])
		c.Locals("user_email", claims["email"])

		if claims["is_admin"] == true {
			c.Locals("is_admin", true)
		} else {
			c.Locals("is_admin", false)
		}
		return c.Next()
	}

	return c.Status(http.StatusUnauthorized).JSON(
		map[string]string{
			"message": "Unauthorized, access token is invalid!",
		},
	)
}

func IsAdmin(c *fiber.Ctx) error {
	if c.Locals("is_admin") == true {
		return c.Next()
	}

	return c.Status(http.StatusUnauthorized).JSON(
		map[string]string{
			"message": "Unauthorized to access this menu",
		},
	)

}
