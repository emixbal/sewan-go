package models

import (
	"log"
	"net/http"
	"sewan-go/config"
	"time"
)

func TransactionDetailAndPayments(transaction_id int) (Response, error) {
	type TransactionDetailAndPaymentsResponse struct {
		Detail interface{} `json:"detail"`
	}
	type Detail struct {
		Pemesan   string    `json:"pemesan"`
		Address   string    `json:"address"`
		Phone     string    `json:"phone"`
		Email     string    `json:"email"`
		StartDate time.Time `json:"start_date"`
		ENdDate   time.Time `json:"end_date"`
	}

	var tdapr TransactionDetailAndPaymentsResponse
	var detail Detail
	var res Response

	db := config.GetDBInstance()

	if r := db.Table("transactions t").
		Select("c.name AS pemesan, c.address, c.phone, c.email, t.start_date, t.end_date").
		Joins("left join customers c ON t.customer_id=c.id").
		Where("t.id = ?", transaction_id).Scan(&detail); r.Error != nil {

		log.Println("TransactionDetailAndPayments")
		log.Println(r.Error)
		res.Status = http.StatusInternalServerError
		res.Message = "err"
		return res, r.Error
	}

	tdapr.Detail = detail

	res.Status = http.StatusInternalServerError
	res.Message = "err"
	res.Data = tdapr

	return res, nil
}
