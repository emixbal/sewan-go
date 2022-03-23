package models

import (
	"errors"
	"fmt"
	"net/http"
	"sewan-go/config"
	"time"

	"gorm.io/gorm"
)

type TransactionItem struct {
	ID            uint `json:"id"`
	TransactionID int  `json:"transaction_id" gorm:"not null"`
	Transaction   Transaction
	ProductID     int `gorm:"not null" json:"product_id"`
	Product       Product
	Qty           int       `gorm:"not null" json:"qty"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP" faker:"-"`
	UpdatedAt     time.Time `json:"update_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" faker:"-"`
}

func AddItemToTransaction(item *TransactionItem) (Response, error) {
	var res Response
	var transaction Transaction

	db := config.GetDBInstance()

	result := db.First(&transaction, item.TransactionID)
	if result.Error != nil {
		if is_notfound := errors.Is(result.Error, gorm.ErrRecordNotFound); is_notfound {
			res.Status = http.StatusOK
			res.Message = "can't find transaction"
			return res, result.Error
		}
	}

	start := transaction.StartDate
	end := transaction.EndDate
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		fmt.Println(d)
	}

	// if result := db.Create(&item); result.Error != nil {
	// 	log.Println("error CreateATransactionItem")
	// 	log.Println("result.Error==>", result.Error)

	// 	res.Status = http.StatusInternalServerError
	// 	res.Message = "error save new record"
	// 	return res, result.Error
	// }

	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = transaction

	return res, nil
}
