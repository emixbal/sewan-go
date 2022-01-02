package controllers

import (
	"fmt"
	"log"
	"net/http"
	"sejuta-cita/app/helpers"
	"sejuta-cita/app/models"
	"sejuta-cita/app/requests"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func UserRegister(c *fiber.Ctx) error {
	var user models.User

	p := new(requests.RegisterForm)
	if err := c.BodyParser(p); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Something went wrong",
		})
	}
	v := validate.Struct(p)
	if !v.Validate() {
		return c.JSON(fiber.Map{
			"message": v.Errors.One(),
		})
	}

	hashPassword, err := helpers.GeneratePassword(p.Password)
	if err != nil {
		fmt.Println(err)
	}

	user.Email = p.Email
	user.Name = p.Name
	user.Password = hashPassword

	result, err := models.UserRegister(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Empty payloads",
		})
	}
	return c.Status(result.Status).JSON(result)
}

func UserLogin(c *fiber.Ctx) error {
	p := new(requests.LoginForm)
	if err := c.BodyParser(p); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Empty payloads",
		})
	}
	v := validate.Struct(p)
	if !v.Validate() {
		return c.JSON(fiber.Map{
			"message": v.Errors.One(),
		})
	}

	isExist, isMatch, tokenString, user, err := models.CheckLogin(p.Email, p.Password)
	if !isExist {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "user not registered",
		})
	}
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "something went wrong",
		})
	}
	if !isMatch {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "password is incorrect",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"token": tokenString,
		"user": fiber.Map{
			"id":       user.ID,
			"email":    user.Email,
			"is_admin": user.IsAdmin,
		},
	})
}

func UserRefreshToken(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(map[string]string{"message": "User Refresh Token"})
}
