package models

import (
	"fmt"
	"net/http"
	"sejuta-cita/config"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	IsAdmin  bool   `json:"is_admin,omitempty" gorm:"default:false"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func FethAllUsers() (Response, error) {
	var books []Book
	var res Response

	db := config.GetDBInstance()

	if result := db.Find(&books); result.Error != nil {
		fmt.Print("error FethAllBooks")
		fmt.Print(result.Error)

		res.Status = http.StatusInternalServerError
		res.Message = "error fetchin records"
		return res, result.Error
	}

	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = books

	return res, nil
}

func CreateAUser(book *User) (Response, error) {
	var res Response
	db := config.GetDBInstance()

	if result := db.Create(&book); result.Error != nil {
		fmt.Print("error CreateABook")
		fmt.Print(result.Error)

		res.Status = http.StatusInternalServerError
		res.Message = "error save new record"
		return res, result.Error
	}

	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = book

	return res, nil
}
