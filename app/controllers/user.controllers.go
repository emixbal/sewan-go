package controllers

import (
	"fmt"
	"log"
	"net/http"
	"sejuta-cita/app/helpers"
	"sejuta-cita/app/models"
	"sejuta-cita/app/requests"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
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
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "invalid user id",
		})
	}
	result, _ := models.ShowUserDetail(user_id)
	return c.Status(result.Status).JSON(result)
}

func UserSoftDelete(c *fiber.Ctx) error {
	result, _ := models.UserSoftDelete(c.Params("id"))

	return c.Status(200).JSON(result)
}

func UserHardDelete(c *fiber.Ctx) error {
	result, _ := models.UserHardDelete(c.Params("id"))

	return c.Status(200).JSON(result)
}

func UserUpdate(c *fiber.Ctx) error {
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

	result, err := models.UserUpdate(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Something went wrong",
		})
	}
	return c.Status(result.Status).JSON(result)
}
