package controllers

import (
	"log"
	"net/http"
	"sewan-go/app/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func TransactionDetailAndPayments(c *fiber.Ctx) error {
	transaction_id, err_transaction_id := strconv.Atoi(c.Params("transaction_id"))
	if err_transaction_id != nil {
		log.Println("err_transaction_id==>", err_transaction_id)
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "transaction id is empty or invalid format"})
	}

	result, _ := models.TransactionDetailAndPayments(transaction_id)
	return c.Status(result.Status).JSON(result)

}
