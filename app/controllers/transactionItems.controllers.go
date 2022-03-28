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

	product_id, err_product_id := strconv.Atoi(c.FormValue("product_id"))
	if err_product_id != nil {
		log.Println("err_product_id==> ", err_product_id)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "product_id is empty or invalid format"})
	}

	qty, err_qty := strconv.Atoi(c.FormValue("qty"))
	if err_qty != nil {
		log.Println("err_qty==> ", err_qty)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "qty is empty or invalid format"})
	}

	item.TransactionID = transaction_id
	item.ProductID = product_id
	item.Qty = qty

	result, _ := models.AddItemToTransaction(&item)
	return c.Status(result.Status).JSON(result)
}

func TransactionItemDelete(c *fiber.Ctx) error {
	result, _ := models.TransactionItemDelete(c.Params("item_id"))

	return c.Status(result.Status).JSON(result)
}
