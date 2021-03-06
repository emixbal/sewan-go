package controllers

import (
	"log"
	"net/http"
	"sewan-go/app/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func AddItemToTransaction(c *fiber.Ctx) error {
	var item models.TransactionItem

	transaction_id, err_transaction_id := strconv.Atoi(c.Params("transaction_id"))
	if err_transaction_id != nil {
		log.Println("err_transaction_id==> ", err_transaction_id)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "transaction_id is empty or invalid format"})
	}

	payload := struct {
		ProductId int `json:"product_id"`
		Qty       int `json:"qty"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	if payload.ProductId == 0 {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "product_id is required"})
	}

	if payload.Qty == 0 {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "qty is required"})
	}

	item.TransactionID = transaction_id
	item.ProductID = payload.ProductId
	item.Qty = payload.Qty

	result, _ := models.AddItemToTransaction(&item)
	return c.Status(result.Status).JSON(result)
}

func TransactionItemDelete(c *fiber.Ctx) error {
	result, _ := models.TransactionItemDelete(c.Params("item_id"))

	return c.Status(result.Status).JSON(result)
}

func TransactionItemUpdateQty(c *fiber.Ctx) error {
	var item models.TransactionItem

	transaction_id := c.Params("transaction_id")

	payload := struct {
		ItemId uint `json:"item_id"`
		Qty    int  `json:"qty"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	item.ID = payload.ItemId
	item.Qty = payload.Qty
	if item.Qty < 1 {
		log.Println("qty is required")
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "qty is required"})
	}
	result, _ := models.TransactionItemUpdateQty(transaction_id, &item)
	return c.Status(result.Status).JSON(result)
}
