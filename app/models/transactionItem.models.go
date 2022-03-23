package models

import (
	"log"
	"net/http"
	"sewan-go/config"
	"time"
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
