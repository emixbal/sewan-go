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

	var sisa int
	start := transaction.StartDate
	end := transaction.EndDate
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		db.Raw(
			`
			SELECT
				(SELECT products.qty FROM products WHERE products.id = ? ) - IFNULL(SUM( ( SELECT transaction_items.qty FROM transaction_items WHERE transaction_items.transaction_id = t.id AND transaction_items.product_id = ? ) ),0) AS sisa
			FROM
				transactions t 
			WHERE
				t.start_date <= ? AND ? < t.end_date
			`,
			item.ProductID, item.ProductID, d, d,
		).Scan(&sisa)

		fmt.Println("sisa===>", sisa)

		if item.Qty > sisa {
			res.Status = http.StatusOK
			res.Message = "kurang"
			return res, nil
		}
	}

	if result := db.Create(&item); result.Error != nil {
		log.Println("error AddItemToTransaction")
		log.Println("result.Error==>", result.Error)

		res.Status = http.StatusInternalServerError
		res.Message = "error save new record"
		return res, result.Error
	}

	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = transaction

	return res, nil
}
