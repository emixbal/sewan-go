package models

import (
	"fmt"
	"net/http"
	"sewan-go/config"
	"time"
)

type Product struct {
	ID        uint      `json:"id"`
	IsActive  bool      `json:"is_active,omitempty" gorm:"default:true"`
	Name      string    `json:"name"`
	Kode      string    `json:"kode"`
	CreatedAt time.Time `json:"created_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP" faker:"-"`
	UpdatedAt time.Time `json:"update_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" faker:"-"`
}

func FethAllProducts(limit, offset int) (Response, error) {
	var products []Product
	var res Response

	db := config.GetDBInstance()

	if result := db.Limit(limit).Offset(offset).Where("is_active = ?", true).Find(&products); result.Error != nil {
		fmt.Print("error FethAllproducts")
		fmt.Print(result.Error)

		res.Status = http.StatusInternalServerError
		res.Message = "error fetchin records"
		return res, result.Error
	}

	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = products

	return res, nil
}

func CreateAProduct(product *Product) (Response, error) {
	var res Response
	db := config.GetDBInstance()

	if result := db.Create(&product); result.Error != nil {
		fmt.Print("error CreateAProduct")
		fmt.Print(result.Error)

		res.Status = http.StatusInternalServerError
		res.Message = "error save new record"
		return res, result.Error
	}

	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = product

	return res, nil
}
