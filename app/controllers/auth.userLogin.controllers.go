package controllers

import (
	"log"
	"net/http"
	"sejuta-cita/app/models"
	"sejuta-cita/app/requests"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

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
