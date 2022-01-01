package models

import (
	"errors"
	"fmt"
	"net/http"
	"sejuta-cita/config"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	IsAdmin  bool   `json:"is_admin,omitempty" gorm:"default:false"`
	Email    string `json:"email" gorm:"index:idx_name,unique"`
	Password string `json:"-"`
}

func FethAllUsers() (Response, error) {
	var users []User
	var res Response

	db := config.GetDBInstance()

	if result := db.Find(&users); result.Error != nil {
		fmt.Print("error FethAllUsers")
		fmt.Print(result.Error)

		res.Status = http.StatusInternalServerError
		res.Message = "error fetchin records"
		return res, result.Error
	}

	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = users

	return res, nil
}

func ShowUserDetail(user_id string) (Response, error) {
	var res Response
	var user User
	db := config.GetDBInstance()

	result := db.First(&user, user_id)
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
	res.Data = user

	return res, nil
}

func CreateAUser(user *User) (Response, error) {
	var res Response
	db := config.GetDBInstance()

	if result := db.Create(&user); result.Error != nil {
		fmt.Print("error CreateAUser")
		fmt.Print(result.Error)

		res.Status = http.StatusInternalServerError
		res.Message = "error save new record"
		return res, result.Error
	}

	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = user

	return res, nil
}
