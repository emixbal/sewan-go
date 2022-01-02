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

	limit := 15
	offset := 0
	if c.Query("per_page") != "" {
		limit, _ = strconv.Atoi(c.Query("per_page"))
		if limit > 50 {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "too much data to show",
			})

		}
	}
	if c.Query("per_page") != "" {
		offset, _ = strconv.Atoi(c.Query("page"))
		if offset != 0 {
			offset = offset - 1
		}
	}

	result, _ := models.FethAllUsers(limit, offset)
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

	return c.Status(result.Status).JSON(result)
}

func UserHardDelete(c *fiber.Ctx) error {
	result, _ := models.UserHardDelete(c.Params("id"))

	return c.Status(result.Status).JSON(result)
}

func UserUpdate(c *fiber.Ctx) error {
	var user models.User

	p := new(requests.UserUpdateForm)
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

	user.Email = p.Email
	user.Name = p.Name

	result, err := models.UserUpdate(&user, c.Params("id"))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(result)
	}
	return c.Status(result.Status).JSON(result)
}

func NewPasswordSelf(c *fiber.Ctx) error {
	user_id_interface := c.Locals("user_id")
	user_id := fmt.Sprintf("%v", user_id_interface)

	p := new(requests.UserUpdatePasswordForm)
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

	hashPassword, err_generate_pass := helpers.GeneratePassword(p.Password)
	if err_generate_pass != nil {
		log.Panicln(err_generate_pass)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "something went wrong",
		})
	}

	result, _ := models.NewPassword(user_id, hashPassword)
	return c.Status(result.Status).JSON(result)
}

func NewPassword(c *fiber.Ctx) error {
	user_id := c.Params("id")

	p := new(requests.UserUpdatePasswordForm)
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

	hashPassword, err_generate_pass := helpers.GeneratePassword(p.Password)
	if err_generate_pass != nil {
		log.Panicln(err_generate_pass)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "something went wrong",
		})
	}

	result, _ := models.NewPassword(user_id, hashPassword)
	return c.Status(result.Status).JSON(result)
}
