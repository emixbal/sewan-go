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

	var sisa int
	start := transaction.StartDate
	end := transaction.EndDate
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {

		fmt.Println("date===>", d)
		qry := fmt.Sprintf("SELECT p.qty-SUM(ti.qty) AS sisa FROM  transaction_items ti LEFT JOIN transactions t ON ti.transaction_id=t.id LEFT JOIN products p ON ti.product_id=p.id WHERE ti.product_id=%d AND (t.start_date <=%d AND %d < t.end_date)", item.ProductID, d, d)
		db.Raw(qry).Scan(&sisa)
		if sisa < item.Qty {
			fmt.Println("is masih? habis")
		} else {
			fmt.Println("is masih? masih")
		}
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
