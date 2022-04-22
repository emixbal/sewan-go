package models

import (
	"log"
	"net/http"
	"sewan-go/config"
	"time"
)

func TransactionDetailAndPayments(transaction_id int) (Response, error) {
	type TransactionDetailAndPaymentsResponse struct {
		Detail           interface{} `json:"detail"`
		PaymentSummary   interface{} `json:"payment_summary"`
		TransactionItems interface{} `json:"transaction_items"`
		PaymentList      interface{} `json:"payment_list"`
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
	type PaymentSummary struct {
		TotalTagihan int `json:"total_tagihan"`
		TotalDibayar int `json:"total_dibayar"`
		SisaTagihan  int `json:"sisa_tagihan"`
	}
	type TransactionItem struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Qty      int    `json:"qty"`
		Price    int    `json:"price"`
		SubTotal int    `json:"sub_total"`
	}

	var tdapr TransactionDetailAndPaymentsResponse
	var detail Detail
	var payment_summary PaymentSummary
	var transactionItems []TransactionItem
	var payments []Payment
	var res Response

	db := config.GetDBInstance()

	if r := db.Table("transactions t").
		Select(`
			t.id,
			c.name AS pemesan, c.address, c.phone, c.email, t.start_date, t.end_date,
			(SELECT IF(st.name IS NULL or st.name = '', 'Default', st.name)) AS status_transaction
		`).
		Joins("left join customers c ON t.customer_id=c.id").
		Joins("left join status_transactions st ON st.id=t.status_transaction_id").
		Where("t.id = ?", transaction_id).Scan(&detail); r.Error != nil {

		log.Println("TransactionDetailAndPayments detail")
		log.Println(r.Error)
		res.Status = http.StatusInternalServerError
		res.Message = "err"
		return res, r.Error
	}

	if r := db.Raw(`SELECT
			t.id AS transaksi_id,
			(SELECT SUM(p.price*ti.qty) FROM transaction_items ti JOIN products p WHERE ti.product_id=p.id AND ti.transaction_id=t.id) AS total_tagihan,
			(SELECT SUM(py.nominal) FROM payments py WHERE py.transaction_id=t.id) AS total_dibayar,
			(SELECT (total_tagihan-total_dibayar)) AS sisa_tagihan
		FROM transactions t
		WHERE t.is_active= ?
		AND t.id= ?
		LIMIT 1`, true, transaction_id).Scan(&payment_summary); r.Error != nil {

		log.Println("TransactionDetailAndPayments payment_summary")
		log.Println(r.Error)
		res.Status = http.StatusInternalServerError
		res.Message = "err"
		return res, r.Error
	}

	if r := db.Table("transaction_items ti").
		Select("ti.id, p.name, ti.qty, p.price, (ti.qty*p.price) as sub_total").
		Joins("left join products p on p.id = ti.product_id").
		Where("ti.transaction_id = ?", transaction_id).
		Scan(&transactionItems); r.Error != nil {

		log.Println("TransactionDetailAndPayments transactionItems")
		log.Println(r.Error)
		res.Status = http.StatusInternalServerError
		res.Message = "err"
		return res, r.Error
	}
	if r := db.Where("transaction_id = ?", transaction_id).Find(&payments); r.Error != nil {
		log.Println("TransactionDetailAndPayments payments")
		log.Println(r.Error)
		res.Status = http.StatusInternalServerError
		res.Message = "err"
		return res, r.Error
	}

	tdapr.Detail = detail
	tdapr.PaymentSummary = payment_summary
	tdapr.TransactionItems = transactionItems
	tdapr.PaymentList = payments

	res.Status = http.StatusOK
	res.Message = config.SuccessMessage
	res.Data = tdapr

	return res, nil
}
