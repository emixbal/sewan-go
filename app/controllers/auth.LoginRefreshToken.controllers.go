package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sewan-go/app/models"
	"sewan-go/app/requests"
	"sewan-go/config"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/golang-jwt/jwt"
	"github.com/gookit/validate"
	"golang.org/x/crypto/bcrypt"
)

var refreshSecret = []byte(os.Getenv("REFRESH_SECRET"))

func LoginRefrehToken(c *fiber.Ctx) error {
	var userClaim models.UserClaim

	p := new(requests.LoginForm)
	if err := c.BodyParser(p); err != nil {
		return err
	}
	v := validate.Struct(p)
	if !v.Validate() {
		return c.JSON(fiber.Map{
			"message": v.Errors.One(),
		})
	}
	u := new(models.User)

	db := config.GetDBInstance()

	if res := db.Where("email = ?", p.Email).First(&u); res.RowsAffected <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid Email!",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p.Password)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Password is incorrect!",
		})
	}
	userClaim.Issuer = utils.UUIDv4()
	userClaim.Id = int(u.ID)
	userClaim.Email = u.Email
	userClaim.IsAdmin = u.IsAdmin
	accessToken, refreshToken := models.GenerateTokens(&userClaim, false)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": fiber.Map{
			"id":       u.ID,
			"email":    u.Email,
			"is_admin": u.IsAdmin,
		},
	})
}

func RefreshToken(c *fiber.Ctx) error {
	var userClaim models.UserClaim

	refreshToken := c.FormValue("refresh_token")

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return refreshSecret, nil
	})

	if err != nil {
		fmt.Println("the error from parse: ", err)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": err})
	}

	//is token valid?
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "StatusUnauthorized",
		})
	}

	user_id := claims["user_id"]
	user_id_int := int(user_id.(float64))
	userClaim.Issuer = fmt.Sprintf("%v", claims["issuer"])
	userClaim.Id = user_id_int
	userClaim.Email = fmt.Sprintf("%v", claims["email"])
	if claims["is_admin"] == true {
		userClaim.IsAdmin = true
	} else {
		userClaim.IsAdmin = false
	}

	// if fail refresh token
	accessToken, refreshToken := models.GenerateTokens(&userClaim, true)
	if len(accessToken) < 1 || len(refreshToken) < 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "something went wrong",
		})
	}

	// set token to blacklist in redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// rdb := config.GetDBInstanceRedis()
	errrdb := rdb.Set(context.TODO(), userClaim.Issuer, userClaim.Id, 0).Err()
	if errrdb != nil {
		log.Println("====>redis err save blacklist token<===")
		log.Println(errrdb)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": fiber.Map{
			"id":       userClaim.Id,
			"email":    userClaim.Email,
			"is_admin": userClaim.IsAdmin,
		},
	})
}
