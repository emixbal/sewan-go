package controllers

import (
	"log"
	"net/http"
	"sewan-go/app/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func NewTransaction(c *fiber.Ctx) error {
	var transaction models.Transaction

	layout := "2006-01-02"
	start_date, err_start := time.Parse(layout, c.FormValue("start_date"))
	if err_start != nil {
		log.Println("err_start==> ", err_start)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "start_date is empty or invalid format"})
	}

	end_date, err_end := time.Parse(layout, c.FormValue("end_date"))
	if err_end != nil {
		log.Println("err_end==> ", err_end)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "end_date is empty or invalid format"})
	}

	customer_id, err_customer_id := strconv.Atoi(c.FormValue("customer_id"))
	if err_customer_id != nil {
		log.Println("err_customer_id==> ", err_customer_id)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "customer_id is empty or invalid format"})
	}

	transaction.StartDate = start_date
	transaction.EndDate = end_date
	transaction.CustomerID = customer_id

	result, _ := models.NewTransaction(&transaction)
	return c.Status(result.Status).JSON(result)
}

func TransactionDetail(c *fiber.Ctx) error {
	transaction_id, err_transaction_id := strconv.Atoi(c.Params("transaction_id"))
	if err_transaction_id != nil {
		log.Println(err_transaction_id)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "transaction_id is empty or invalid format"})
	}
	result, _ := models.TransactionDetail(transaction_id)
	return c.Status(result.Status).JSON(result)
}

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

	result, _ := models.NewTransactionItem(&item)
	return c.Status(result.Status).JSON(result)
}
