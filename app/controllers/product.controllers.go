package controllers

import (
	"net/http"
	"sewan-go/app/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func FetchAllproducts(c *fiber.Ctx) error {

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

	result, _ := models.FethAllProducts(limit, offset)
	return c.Status(result.Status).JSON(result)
}

func CreateANewProduct(c *fiber.Ctx) error {
	var product models.Product

	product.Name = c.FormValue("name")
	product.Kode = c.FormValue("kode")

	if product.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "name is required"})
	}
	if product.Kode == "" {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "no_isbn is required"})
	}

	result, _ := models.CreateAProduct(&product)
	return c.Status(result.Status).JSON(result)
}
