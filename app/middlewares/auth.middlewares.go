package middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var (
	jwtKey = []byte(os.Getenv("JWT_SECRET"))
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
		return jwtKey, nil
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

		// set token to blacklist in redis
		// rdb := config.GetDBInstanceRedis()
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		issuer := fmt.Sprintf("%v", claims["issuer"])
		var ctx = context.TODO()
		val, errrdb := rdb.Get(ctx, issuer).Result()
		if errrdb != nil {
			log.Println("====>redis err read blacklist token<===")
			log.Println(errrdb)
		}

		fmt.Println("val>>>>>", val)

		c.Locals("user_id", claims["user_id"])
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
