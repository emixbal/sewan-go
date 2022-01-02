package controllers

import (
	"fmt"
	"net/http"
	"os"
	"sejuta-cita/app/models"
	"sejuta-cita/app/requests"
	"sejuta-cita/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/gookit/validate"
	"golang.org/x/crypto/bcrypt"
)

var jwtRefreshKey = []byte(os.Getenv("REFRESH_SECRET"))

func LoginRefrehToken(c *fiber.Ctx) error {
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

	accessToken, refreshToken := models.GenerateTokens(u.Email)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.FormValue("refresh_token")

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtRefreshKey, nil
	})

	if err != nil {
		fmt.Println("the error from parse: ", err)
		return c.Status(500).JSON(fiber.Map{"message": err})
	}

	//is token valid?
	_, ok := token.Claims.(jwt.Claims)

	if !ok {
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"message": "StatusUnauthorized",
		})
	}

	claims, _ := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims

	fmt.Println(claims)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "StatusUnauthorized",
	})
}
