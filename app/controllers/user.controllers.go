package controllers

import (
	"log"
	"net/http"
	"sejuta-cita/app/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func FetchAllUsers(c *fiber.Ctx) error {
	result, _ := models.FethAllUsers()
	return c.Status(result.Status).JSON(result)
}

func ShowUserDetail(c *fiber.Ctx) error {
	user_id := c.Params("user_id")
	_, err := strconv.Atoi(user_id)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "invalid user id"})
	}
	result, _ := models.ShowUserDetail(user_id)
	return c.Status(result.Status).JSON(result)
}

func CreateUser(c *fiber.Ctx) error {
	var user models.User

	user.Name = c.FormValue("name")
	user.Email = c.FormValue("email")

	if user.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "name is required"})
	}
	if user.Email == "" {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "no_isbn is required"})
	}

	result, _ := models.CreateAUser(&user)
	return c.Status(result.Status).JSON(result)
}
