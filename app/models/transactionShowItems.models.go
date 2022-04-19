package models

import (
	"log"
	"net/http"
	"sewan-go/config"
)

func TransactionShowItems(id int) (Response, error) {
	type TransactionItemsDetailResponse struct {
		TransactionItems interface{} `json:"transaction_items"`
		Payment          interface{} `json:"payment"`
	}

	type TransactionItemRes struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Qty      int    `json:"qty"`
		Price    int    `json:"price"`
		SubTotal int    `json:"sub_total"`
	}

	var res Response
	var tidr TransactionItemsDetailResponse
	var payment_summary PaymentSummary
	var arrTransactionItem []TransactionItemRes
	var arrTransactionItemRes []TransactionItemRes
	db := config.GetDBInstance()

	result := db.Table("transaction_items ti").
		Select("ti.id, p.name, ti.qty, p.price").
		Joins("left join products p on p.id = ti.product_id").
		Where("ti.transaction_id = ?", id).
		Scan(&arrTransactionItem)
	if result.Error != nil {
		log.Println("err TransactionShowItems")
		log.Println(result.Error)
		res.Status = http.StatusInternalServerError
		res.Message = "Something went wrong!"
		return res, nil
	}

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
		log.Println("err TransactionShowItems res_qry")
		log.Println(res_qry.Error)
		res.Status = http.StatusInternalServerError
		res.Message = "Something went wrong!"
		return res, nil
	}

	for _, item := range arrTransactionItem {
		item.SubTotal = item.Price * item.Qty
		arrTransactionItemRes = append(arrTransactionItemRes, item)
	}

	tidr.TransactionItems = arrTransactionItemRes
	tidr.Payment = payment_summary

	res.Status = http.StatusOK
	res.Message = config.SuccessMessage
	res.Data = tidr

	return res, nil
}
