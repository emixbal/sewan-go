package controllers

import (
	"log"
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

	if c.Query("page") != "" {
		offset, _ = strconv.Atoi(c.Query("page"))
		if offset != 0 {
			offset = offset - 1
		}
	}

	result, _ := models.FethAllProducts(limit, offset)
	return c.Status(result.Status).JSON(result)
}

func ShowProductDetail(c *fiber.Ctx) error {
	product_id, err_product_id := strconv.Atoi(c.Params("product_id"))
	if err_product_id != nil {
		log.Println("err_product_id==>", err_product_id)
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "product id is empty or invalid format"})
	}
	result, _ := models.DetailProduct(product_id)
	return c.Status(result.Status).JSON(result)
}

func CreateANewProduct(c *fiber.Ctx) error {
	var product models.Product

	payload := struct {
		Name string `json:"name"`
		Kode string `json:"kode"`
		Qty  int    `json:"qty"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	product.Name = payload.Name
	product.Kode = payload.Kode
	product.Qty = payload.Qty

	if product.Name == "" {
		log.Println("name is required")
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "name is required"})
	}
	if product.Kode == "" {
		log.Println("kode is required")
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "kode is required"})
	}
	if product.Qty < 1 {
		log.Println("qty is required")
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "qty is required"})
	}

	result, _ := models.CreateAProduct(&product)
	return c.Status(result.Status).JSON(result)
}

func UpdateProduct(c *fiber.Ctx) error {
	var product models.Product

	product_id := c.Params("product_id")

	payload := struct {
		Name string `json:"name"`
		Kode string `json:"kode"`
		Qty  int    `json:"qty"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	product.Name = payload.Name
	product.Kode = payload.Kode
	product.Qty = payload.Qty

	if product_id == "" {
		log.Println("product_id is required")
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "product_id is required"})
	}
	if product.Name == "" {
		log.Println("name is required")
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "name is required"})
	}
	if product.Kode == "" {
		log.Println("kode is required")
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "kode is required"})
	}
	if product.Qty <= 1 {
		log.Println("qty is required")
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "qty is required"})
	}

	result, _ := models.UpdateProduct(&product, product_id)
	return c.Status(result.Status).JSON(result)
}

func ProductSoftDelete(c *fiber.Ctx) error {
	result, _ := models.ProductSoftDelete(c.Params("product_id"))

	return c.Status(result.Status).JSON(result)
}
