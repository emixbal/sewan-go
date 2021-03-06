package controllers

import (
	"net/http"
	"sewan-go/app/models"

	"github.com/gofiber/fiber/v2"
)

func FetchAllBooks(c *fiber.Ctx) error {
	result, _ := models.FethAllBooks()
	return c.Status(result.Status).JSON(result)
}

func CreateBook(c *fiber.Ctx) error {
	var book models.Book

	book.Author = c.FormValue("author")
	book.Name = c.FormValue("name")
	book.NoISBN = c.FormValue("no_isbn")

	if book.Author == "" {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "author is required"})
	}
	if book.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "name is required"})
	}
	if book.NoISBN == "" {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "no_isbn is required"})
	}

	result, _ := models.CreateABook(&book)
	return c.Status(result.Status).JSON(result)
}
