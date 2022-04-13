package models

import (
	"errors"
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

	var total_prod_in_transaction int64
	prod_is_exist := db.Model(&TransactionItem{}).Where("transaction_id = ? AND product_id = ?", item.TransactionID, item.ProductID).Count(&total_prod_in_transaction)
	if prod_is_exist.Error != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "something went wrong"
		return res, prod_is_exist.Error
	}

	if total_prod_in_transaction > 0 {
		log.Println("=================PRODUCT ALREADY EXIST=================")
		log.Println("total_prod_in_transaction==>", total_prod_in_transaction)

		res.Status = http.StatusBadRequest
		res.Message = "product already exist"
		return res, nil
	}

	var sisa int
	start := transaction.StartDate
	end := transaction.EndDate
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		db.Raw(
			`
			SELECT
				((SELECT products.qty FROM products WHERE products.id = ? ) - IFNULL(SUM((SELECT ti.qty FROM transaction_items ti WHERE ti.transaction_id=t.id AND ti.product_id=?)),0)) AS sisa
			FROM
				transactions t 
			WHERE
				t.start_date <= ?
				AND ? <= t.end_date
			`,
			item.ProductID, item.ProductID, d, d,
		).Scan(&sisa)

		if item.Qty > sisa {
			log.Println("=================OUT OF STOCK=================")
			log.Println("sisa==>", sisa)
			log.Println("item.Qty==>", item.Qty)

			res.Status = http.StatusBadRequest
			res.Message = "out of stock"
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
	res.Message = config.SuccessMessage
	return res, nil
}

func TransactionItemDelete(item_id string) (Response, error) {
	var res Response
	var item TransactionItem

	db := config.GetDBInstance()
	result := db.Unscoped().Delete(&item, item_id)

	if result.Error != nil {
		log.Println(result.Error)
		res.Status = http.StatusInternalServerError
		res.Message = "Error"

		return res, result.Error
	}
	res.Status = http.StatusOK
	res.Message = config.SuccessMessage

	return res, nil
}

func TransactionItemUpdateQty(transaction_id string, item_payload *TransactionItem) (Response, error) {
	var res Response
	var item TransactionItem

	db := config.GetDBInstance()
	result := db.Where("id = ?", item_payload.ID).Take(&item)
	if result.Error != nil {
		if is_notfound := errors.Is(result.Error, gorm.ErrRecordNotFound); is_notfound {
			res.Status = http.StatusOK
			res.Message = "can't find record"
			return res, result.Error
		}

		log.Println("error TransactionItemUpdateQty")
		log.Println(result.Error)

		res.Status = http.StatusInternalServerError
		res.Message = "something went wrong"
		return res, result.Error
	}

	item.Qty = item_payload.Qty

	db.Save(&item)

	res.Status = http.StatusOK
	res.Message = config.SuccessMessage
	res.Data = item

	return res, nil
}
