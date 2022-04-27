package models

import (
	"errors"
	"log"
	"net/http"
	"sewan-go/config"
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID                  uint `json:"id"`
	IsActive            bool `json:"is_active,omitempty" gorm:"default:true"`
	CustomerID          int  `json:"customer_id"`
	Customer            Customer
	TransactionItems    []TransactionItem `json:"transaction_items"`
	StartDate           time.Time         `gorm:"not null" json:"start_date"`
	EndDate             time.Time         `gorm:"not null" json:"end_date"`
	StatusTransaction   StatusTransaction
	StatusTransactionID int       `json:"status_transaction_id"`
	CreatedAt           time.Time `json:"created_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP" faker:"-"`
	UpdatedAt           time.Time `json:"update_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" faker:"-"`
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
		log.Print("error CreateATransaction")
		log.Print(result.Error)

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
		log.Println(res_qry.Error)
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
		Id                int    `json:"id"`
		StartDate         string `json:"start_date"`
		EndDate           string `json:"end_date"`
		Pemesan           string `json:"pemesan"`
		StatusTransaction string `json:"status_transaction"`
		TotalTagihan      uint   `json:"total_tagihan"`
		TotalDibayar      uint   `json:"total_dibayar"`
		SisaTagihan       uint   `json:"sisa_tagihan"`
	}

	var transactionRes []TransactionRes
	var res Response

	db := config.GetDBInstance()

	result := db.Raw(`
		SELECT
			t.id, t.start_date, t.end_date, c.name AS pemesan,
			(SELECT st.name FROM status_transactions st WHERE st.id = t.status_transaction_id) AS status_transaction,
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
	type TransactionAddPaymentResponse struct {
		Detail         interface{} `json:"transaction_detail"`
		PaymentSummary interface{} `json:"payment_sumamry"`
	}

	type Detail struct {
		ID                string    `json:"id"`
		Pemesan           string    `json:"pemesan"`
		Address           string    `json:"address"`
		Phone             string    `json:"phone"`
		Email             string    `json:"email"`
		StartDate         time.Time `json:"start_date"`
		EndDate           time.Time `json:"end_date"`
		StatusTransaction string    `json:"status_transaction"`
	}

	var tapr TransactionAddPaymentResponse
	var detail Detail
	var payment_summary PaymentSummary
	var res Response
	db := config.GetDBInstance()

	if res_qry := db.Raw(`SELECT
			t.id AS transaksi_id,
			(SELECT SUM(p.price*ti.qty) FROM transaction_items ti JOIN products p WHERE ti.product_id=p.id AND ti.transaction_id=t.id) AS total_tagihan,
			(SELECT IF(SUM(py.nominal) IS NULL OR SUM(py.nominal)='', 0, SUM(py.nominal)) FROM payments py WHERE py.transaction_id=t.id) AS total_dibayar,
			(SELECT (total_tagihan-total_dibayar)) AS sisa_tagihan
		FROM transactions t
		WHERE t.is_active= ?
		AND t.id= ?
		LIMIT 1`, true, payment.TransactionID).Scan(&payment_summary); res_qry.Error != nil {
		log.Println(res_qry.Error)
		res.Status = http.StatusInternalServerError
		res.Message = "err"

		return res, nil
	}

	if payment.Nominal > payment_summary.SisaTagihan {
		log.Println("payment.Nominal", payment.Nominal)
		log.Println("payment_summary.SisaTagihan", payment_summary.SisaTagihan)
		log.Println("payment.Nominal > payment_summary.SisaTagihan")
		res.Status = http.StatusBadRequest
		res.Message = "Terlalu banyak"

		return res, nil
	}

	if result := db.Create(&payment); result.Error != nil {
		log.Print("error TransactionAddPayment")
		log.Print(result.Error)

		res.Status = http.StatusInternalServerError
		res.Message = "error save new record"
		return res, result.Error
	}

	if r := db.Table("transactions t").
		Select(`
			t.id,
			c.name AS pemesan, c.address, c.phone, c.email, t.start_date, t.end_date,
			(SELECT IF(st.name IS NULL or st.name = '', 'Default', st.name)) AS status_transaction
		`).
		Joins("left join customers c ON t.customer_id=c.id").
		Joins("left join status_transactions st ON st.id=t.status_transaction_id").
		Where("t.id = ?", payment.TransactionID).Scan(&detail); r.Error != nil {

		log.Println("TransactionDetailAndPayments detail")
		log.Println(r.Error)
		res.Status = http.StatusInternalServerError
		res.Message = "err"
		return res, r.Error
	}

	if res_qry := db.Raw(`SELECT
			t.id AS transaksi_id,
			(SELECT SUM(p.price*ti.qty) FROM transaction_items ti JOIN products p WHERE ti.product_id=p.id AND ti.transaction_id=t.id) AS total_tagihan,
			(SELECT SUM(py.nominal) FROM payments py WHERE py.transaction_id=t.id) AS total_dibayar,
			(SELECT (total_tagihan-total_dibayar)) AS sisa_tagihan
		FROM transactions t
		WHERE t.is_active= ?
		AND t.id= ?
		LIMIT 1`, true, payment.TransactionID).Scan(&payment_summary); res_qry.Error != nil {
		log.Println(res_qry.Error)
		res.Status = http.StatusInternalServerError
		res.Message = "err"

		return res, nil
	}

	tapr.PaymentSummary = payment_summary
	tapr.Detail = detail

	res.Status = http.StatusOK
	res.Message = config.SuccessMessage
	res.Data = tapr

	return res, nil
}

func TransactionChangeStatus(transaction_id int, status int) (Response, error) {
	var res Response
	var transaction Transaction

	db := config.GetDBInstance()
	result := db.Where("id = ?", transaction_id).Take(&transaction)
	if result.Error != nil {
		if is_notfound := errors.Is(result.Error, gorm.ErrRecordNotFound); is_notfound {
			res.Status = http.StatusOK
			res.Message = "can't find record"
			return res, result.Error
		}

		log.Println("TransactionChangeSendStatus find err")
		log.Println(result.Error)

		res.Status = http.StatusInternalServerError
		res.Message = "something went wrong"
		return res, result.Error
	}

	transaction.StatusTransactionID = status

	if r := db.Save(&transaction); r.Error != nil {
		log.Println("TransactionChangeSendStatus update status err")
		log.Println(r.Error)
		res.Status = http.StatusInternalServerError
		res.Message = "error update status"
		return res, result.Error
	}

	res.Status = http.StatusOK
	res.Message = config.SuccessMessage

	return res, nil
}

func TransactionListDemage(transaction_id int) (Response, error) {
	type ListDemagesRes struct {
		Id           int    `json:"id"`
		ProductName  string `json:"product_name"`
		Qty          string `json:"qty"`
		ProductPrice string `json:"product_price"`
		SubTotal     string `json:"sub_total"`
	}

	var res Response
	var lsit_demages []ListDemagesRes

	db := config.GetDBInstance()
	if result := db.Table("demages d").
		Select("d.id, p.name AS product_name, d.qty AS qty, p.price AS product_price, (SELECT(d.qty * p.price)) AS sub_total").
		Joins("left join products p on p.id = d.product_id").
		Where("d.transaction_id = ?", transaction_id).
		Scan(&lsit_demages); result.Error != nil {
		log.Println("err TransactionShowItems")
		log.Println(result.Error)
		res.Status = http.StatusInternalServerError
		res.Message = "Something went wrong!"
		return res, nil
	}

	res.Message = config.SuccessMessage
	res.Status = http.StatusOK
	res.Data = lsit_demages

	return res, nil
}
