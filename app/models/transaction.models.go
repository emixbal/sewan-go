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

type PaymentSummary struct {
	ID           int `json:"id"`
	TotalTagihan int `json:"total_tagihan"`
	TotalDibayar int `json:"total_dibayar"`
	SisaTagihan  int `json:"sisa_tagihan"`
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
	type TransactionDetailResponse struct {
		Transaction interface{} `json:"transaction"`
		Payment     interface{} `json:"payment"`
	}

	var res Response
	var tdr TransactionDetailResponse
	var payment_summary PaymentSummary

	var transaction Transaction
	db := config.GetDBInstance()

	result := db.Where("is_active = ?", true).
		Preload("Customer").
		Preload("TransactionItems").
		Preload("TransactionItems.Product").First(&transaction, id)

	res_qry := db.Raw(`SELECT
			t.id AS transaksi_id,
			(SELECT SUM(p.price*ti.qty) FROM transaction_items ti JOIN products p WHERE ti.product_id=p.id AND ti.transaction_id=t.id) AS total_tagihan,
			(SELECT SUM(py.nominal) FROM payments py WHERE py.transaction_id=t.id) AS total_dibayar,
			(SELECT (total_tagihan-total_dibayar)) AS sisa_tagihan
		FROM transactions t
		WHERE t.is_active= ?
		AND t.id= ?
		LIMIT 1`, true, id).Scan(&payment_summary)
	if res_qry.Error != nil {
		fmt.Println(res_qry.Error)
		res.Status = http.StatusInternalServerError
		res.Message = "err"

		return res, nil
	}

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
	tdr.Transaction = transaction
	tdr.Payment = payment_summary

	res.Status = http.StatusOK
	res.Message = config.SuccessMessage
	res.Data = tdr

	return res, nil
}

func TransactionList(limit, offset int) (Response, error) {
	type TransactionRes struct {
		Id           int    `json:"id"`
		StartDate    string `json:"start_date"`
		EndDate      string `json:"end_date"`
		Pemesan      string `json:"pemesan"`
		TotalTagihan uint   `json:"total_tagihan"`
		TotalDibayar uint   `json:"total_dibayar"`
		SisaTagihan  uint   `json:"sisa_tagihan"`
	}

	var transactionRes []TransactionRes
	var res Response

	db := config.GetDBInstance()

	result := db.Raw(`
		SELECT
			t.id, t.start_date, t.end_date, c.name AS pemesan,
			(SELECT SUM(p.price*ti.qty) FROM transaction_items ti JOIN products p WHERE ti.product_id=p.id AND ti.transaction_id=t.id) AS total_tagihan,
			(SELECT SUM(py.nominal) FROM payments py WHERE py.transaction_id=t.id) AS total_dibayar,
			(SELECT (total_tagihan-total_dibayar)) AS sisa_tagihan
		FROM transactions t
		LEFT JOIN customers c ON c.id = t.customer_id
		WHERE t.is_active=?`, true).Scan(&transactionRes)
	if result.Error != nil {
		log.Println("err TransactionList")
		log.Println(result.Error)

		res.Status = http.StatusInternalServerError
		res.Message = "error fetchin records"
		return res, result.Error
	}

	res.Status = http.StatusOK
	res.Message = config.SuccessMessage
	res.Data = transactionRes

	return res, nil
}

func TransactionAddPayment(payment *Payment) (Response, error) {
	var payment_summary PaymentSummary
	var res Response
	db := config.GetDBInstance()

	res_qry := db.Raw(`SELECT
			t.id AS transaksi_id,
			(SELECT SUM(p.price*ti.qty) FROM transaction_items ti JOIN products p WHERE ti.product_id=p.id AND ti.transaction_id=t.id) AS total_tagihan,
			(SELECT SUM(py.nominal) FROM payments py WHERE py.transaction_id=t.id) AS total_dibayar,
			(SELECT (total_tagihan-total_dibayar)) AS sisa_tagihan
		FROM transactions t
		WHERE t.is_active= ?
		AND t.id= ?
		LIMIT 1`, true, payment.TransactionID).Scan(&payment_summary)
	if res_qry.Error != nil {
		fmt.Println(res_qry.Error)
		res.Status = http.StatusInternalServerError
		res.Message = "err"

		return res, nil
	}

	if payment.Nominal > payment_summary.SisaTagihan {
		fmt.Println("payment.Nominal > payment_summary.SisaTagihan")
		res.Status = http.StatusBadRequest
		res.Message = "Terlalu banyak"

		return res, nil
	}

	if result := db.Create(&payment); result.Error != nil {
		fmt.Print("error TransactionAddPayment")
		fmt.Print(result.Error)

		res.Status = http.StatusInternalServerError
		res.Message = "error save new record"
		return res, result.Error
	}

	res.Status = http.StatusOK
	res.Message = config.SuccessMessage

	return res, nil
}
