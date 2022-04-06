package controllers

import (
	"log"
	"net/http"
	"sewan-go/app/models"
	"sewan-go/config"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TransactionCustomerNew(c *fiber.Ctx) error {
	var customer models.Customer
	var transaction models.Transaction

	payload := struct {
		Name      string `json:"name"`
		Address   string `json:"address"`
		Phone     string `json:"phone"`
		Email     string `json:"email"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	customer.Name = payload.Name
	customer.Address = payload.Address
	customer.Phone = payload.Phone
	customer.Email = payload.Email

	if customer.Name == "" {
		log.Println("name is required")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "name is required"})
	}

	if customer.Address == "" {
		log.Println("address is required")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "address is required"})
	}

	if customer.Phone == "" {
		log.Println("phone is required")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "phone is required"})
	}

	layout := "2006-01-02"
	start_date, err_start := time.Parse(layout, payload.StartDate)
	if err_start != nil {
		log.Println("err_start==> ", err_start)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "start_date is empty or invalid format"})
	}

	end_date, err_end := time.Parse(layout, payload.EndDate)
	if err_end != nil {
		log.Println("err_end==> ", err_end)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "end_date is empty or invalid format"})
	}

	if start_date.Before(end_date) == false {
		log.Println("end_date must be after start_date")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "end_date must be after start_date"})
	}

	// save new customer
	result, _ := models.CustomerNewResultId(&customer)

	if result.Message != config.SuccessMessage {
		return c.Status(result.Status).JSON(result)
	}

	transaction.StartDate = start_date
	transaction.EndDate = end_date
	transaction.CustomerID = int(result.CustomerId) // data from customer id

	result_trx, _ := models.TransactionNew(&transaction)
	return c.Status(result_trx.Status).JSON(result_trx)
}
