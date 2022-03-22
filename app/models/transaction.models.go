package models

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sewan-go/config"
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID               uint `json:"id"`
	IsActive         bool `json:"is_active,omitempty" gorm:"default:true"`
	CustomerID       int  `json:"customer_id"`
	Customer         Customer
	TransactionItems []TransactionItem `json:"transaction_items"`
	StartDate        time.Time         `gorm:"not null" json:"start_date"`
	EndDate          time.Time         `gorm:"not null" json:"end_date"`
	CreatedAt        time.Time         `json:"created_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP" faker:"-"`
	UpdatedAt        time.Time         `json:"update_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" faker:"-"`
}

func NewTransaction(transaction *Transaction) (Response, error) {
	var res Response
	db := config.GetDBInstance()

	if result := db.Create(&transaction); result.Error != nil {
		fmt.Print("error CreateATransaction")
		fmt.Print(result.Error)

		res.Status = http.StatusInternalServerError
		res.Message = "error save new record"
		return res, result.Error
	}

	res.Status = http.StatusOK
	res.Message = "success"

	return res, nil
}

func TransactionDetail(id int) (Response, error) {
	var res Response
	var transaction Transaction
	db := config.GetDBInstance()

	result := db.Where("is_active = ?", true).Preload("TransactionItems").First(&transaction, id)
	if result.Error != nil {
		if is_notfound := errors.Is(result.Error, gorm.ErrRecordNotFound); is_notfound {
			res.Status = http.StatusOK
			res.Message = "can't find record"
			return res, result.Error
		}

		res.Status = http.StatusInternalServerError
		res.Message = "Something went wrong!"
		return res, result.Error
	}
	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = transaction

	return res, nil
}

func AddItemToTransaction(item *TransactionItem) (Response, error) {
	var res Response
	db := config.GetDBInstance()

	if result := db.Create(&item); result.Error != nil {
		log.Println("error CreateATransactionItem")
		log.Println("result.Error==>", result.Error)

		res.Status = http.StatusInternalServerError
		res.Message = "error save new record"
		return res, result.Error
	}

	res.Status = http.StatusOK
	res.Message = "success"

	return res, nil
}
