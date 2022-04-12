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

type TransactionItemRes struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Qty      int    `json:"qty"`
	Price    int    `json:"price"`
	SubTotal int    `json:"sub_total"`
}

func TransactionNew(transaction *Transaction) (Response, error) {
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
	res.Message = config.SuccessMessage
	res.Data = transaction

	return res, nil
}

func TransactionDetail(id int) (Response, error) {
	var res Response
	var transaction Transaction
	db := config.GetDBInstance()

	result := db.Where("is_active = ?", true).Preload("Customer").Preload("TransactionItems").First(&transaction, id)
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
	res.Message = config.SuccessMessage
	res.Data = transaction

	return res, nil
}

func TransactionShowItems(id int) (Response, error) {
	var res Response
	var arrTransactionItemRes []TransactionItemRes
	db := config.GetDBInstance()

	result := db.Table("transaction_items ti").Select("ti.id, p.name, ti.qty, p.price").Joins("left join products p on p.id = ti.product_id").Scan(&arrTransactionItemRes)
	if result.Error != nil {
		log.Println("err TransactionShowItems")
		log.Println(result.Error)
		res.Status = http.StatusInternalServerError
		res.Message = "Something went wrong!"
		return res, nil
	}

	for _, item := range arrTransactionItemRes {
		item.SubTotal = item.Price * item.Qty
		arrTransactionItemRes = append(arrTransactionItemRes, item)
	}

	res.Status = http.StatusOK
	res.Message = config.SuccessMessage
	res.Data = arrTransactionItemRes

	return res, nil
}
