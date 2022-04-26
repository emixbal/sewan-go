package models

import (
	"log"
	"net/http"
	"sewan-go/config"
)

type DemagePayload struct {
	ProductID int `json:"product_id"`
	Qty       int `json:"qty"`
}

func TransactionAddDemage(transaction_id int, payload []DemagePayload) (Response, error) {
	var res Response

	var demage Demage
	var demages []Demage

	for _, data := range payload {
		demage.TransactionID = transaction_id
		demage.ProductID = data.ProductID
		demage.Qty = data.Qty

		demages = append(demages, demage)
	}

	db := config.GetDBInstance()
	if r := db.Create(&demages); r.Error != nil {
		log.Println("TransactionAddDemage err")
		log.Println(r.Error)

		res.Message = "TransactionAddDemage err"
		res.Status = http.StatusInternalServerError
		return res, nil
	}

	res.Message = config.SuccessMessage
	res.Status = http.StatusOK

	return res, nil
}
