package controllers

import (
	"log"
	"net/http"
	"sewan-go/app/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TransactionNew(c *fiber.Ctx) error {
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

	if start_date.Before(end_date) == false {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "end_date must be after start_date"})
	}

	customer_id, err_customer_id := strconv.Atoi(c.FormValue("customer_id"))
	if err_customer_id != nil {
		log.Println("err_customer_id==> ", err_customer_id)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "customer_id is empty or invalid format"})
	}

	transaction.StartDate = start_date
	transaction.EndDate = end_date
	transaction.CustomerID = customer_id

	result, _ := models.TransactionNew(&transaction)
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

func TransactionShowItems(c *fiber.Ctx) error {
	transaction_id, err_transaction_id := strconv.Atoi(c.Params("transaction_id"))
	if err_transaction_id != nil {
		log.Println(err_transaction_id)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "transaction_id is empty or invalid format"})
	}
	result, _ := models.TransactionShowItems(transaction_id)
	return c.Status(result.Status).JSON(result)
}

func TransactionList(c *fiber.Ctx) error {
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

	result, _ := models.TransactionList(limit, offset)
	return c.Status(result.Status).JSON(result)
}
