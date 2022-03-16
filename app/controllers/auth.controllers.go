package controllers

import (
	"fmt"
	"log"
	"net/http"
	"sewan-go/app/helpers"
	"sewan-go/app/models"
	"sewan-go/app/requests"

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

func UserRefreshToken(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(map[string]string{"message": "User Refresh Token"})
}
