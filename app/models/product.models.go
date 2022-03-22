package models

import (
	"errors"
	"fmt"
	"net/http"
	"sewan-go/config"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID        uint      `json:"id"`
	IsActive  bool      `json:"is_active,omitempty" gorm:"default:true"`
	Name      string    `json:"name" gorm:"not null default:''"`
	Kode      string    `json:"kode" gorm:"not null default:''"`
	Qty       int       `json:"qty" gorm:"not null default:1"`
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

func DetailProduct(product_id int) (Response, error) {
	var res Response
	var product Product
	db := config.GetDBInstance()

	result := db.Where("is_active = ?", true).First(&product, product_id)
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
	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = product

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
