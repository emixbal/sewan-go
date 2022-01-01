package controllers

import (
	"fmt"
	"net/http"
	"sejuta-cita/app/helpers"
	"sejuta-cita/app/models"

	"github.com/gofiber/fiber/v2"
)

func UserRegister(c *fiber.Ctx) error {
	var user models.User

	txtPassword := c.FormValue("password")

	hashPassword, err := helpers.GeneratePassword(txtPassword)
	if err != nil {
		fmt.Println(err)
	}

	user.Email = c.FormValue("email")
	user.Name = c.FormValue("name")
	user.Password = hashPassword

	if user.Email == "" {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "email is required"})
	}
	if user.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "name is required"})
	}

	result, err := models.UserRegister(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(map[string]string{"message": "something went wrong!"})
	}
	return c.Status(result.Status).JSON(result)
}

func UserLogin(c *fiber.Ctx) error {
	email := c.FormValue("email")
	txtUnHashPassword := c.FormValue("password")

	if email == "" {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "email tidak boleh kosong"})
	}
	if txtUnHashPassword == "" {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "password tidak boleh kosong"})
	}

	isExist, isMatch, tokenString, err := models.CheckLogin(email, txtUnHashPassword)
	if !isExist {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "user not registered"})
	}
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(map[string]string{"message": "something went wrong!"})
	}
	if !isMatch {
		return c.Status(http.StatusUnauthorized).JSON(map[string]string{"message": "Password salah"})
	}

	return c.Status(http.StatusOK).JSON(map[string]string{
		"token": tokenString,
	})
}

func UserRefreshToken(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(map[string]string{"message": "User Refresh Token"})
}
