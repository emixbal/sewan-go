package models

import (
	"fmt"
	"net/http"
	"sewan-go/config"
	"time"
)

type Customer struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Email     string    `json:"email" gorm:"index:idx_email,unique"`
	Phone     string    `json:"phone" gorm:"index:idx_phone,unique"`
	IsActive  bool      `json:"is_active,omitempty" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP" faker:"-"`
	UpdatedAt time.Time `json:"update_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" faker:"-"`
}

type CustomerResponse struct {
	Status     int    `json:"status"`
	Message    string `json:"message"`
	CustomerId uint   `json:"customer_id"`
}

func CustomerNewResultId(customer *Customer) (CustomerResponse, error) {
	var res CustomerResponse
	db := config.GetDBInstance()

	if result := db.Create(&customer); result.Error != nil {
		fmt.Print("error CustomerNewResultId")
		fmt.Print(result.Error)

		res.Status = http.StatusInternalServerError
		res.Message = "err CustomerNewResultId"
		return res, result.Error
	}

	res.Status = http.StatusOK
	res.Message = config.SuccessMessage
	res.CustomerId = customer.ID

	return res, nil
}
