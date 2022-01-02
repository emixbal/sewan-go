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
	IsActive bool   `json:"is_active,omitempty" gorm:"default:true"`
	Email    string `json:"email" gorm:"index:idx_name,unique"`
	Password string `json:"-"`
}

func FethAllUsers(limit, offset int) (Response, error) {
	var users []User
	var res Response

	db := config.GetDBInstance()

	if result := db.Limit(limit).Offset(offset).Where("is_active = ?", true).Find(&users); result.Error != nil {
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

	result := db.Where("is_active = ?", true).First(&user, user_id)
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

func UserSoftDelete(user_id string) (Response, error) {
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

	if !user.IsActive {
		res.Status = http.StatusOK
		res.Message = "user already inactive"

		return res, nil
	}

	user.IsActive = false

	db.Save(&user)

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = user

	return res, nil
}

func UserUpdate(user_payload *User, user_id string) (Response, error) {
	var res Response
	var user User

	db := config.GetDBInstance()
	result := db.Where("id = ?", user_id).Take(&user)
	if result.Error != nil {
		if is_notfound := errors.Is(result.Error, gorm.ErrRecordNotFound); is_notfound {
			res.Status = http.StatusOK
			res.Message = "can't find record"
			return res, result.Error
		}

		res.Status = http.StatusInternalServerError
		res.Message = "something went wrong"
		return res, result.Error
	}

	user.Email = user_payload.Email
	user.Name = user_payload.Name

	db.Save(&user)

	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = user

	return res, nil
}

func UserHardDelete(user_id string) (Response, error) {
	var res Response
	var user User

	db := config.GetDBInstance()
	result := db.Unscoped().Delete(&user, user_id)

	if result.Error != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "Error"

		return res, result.Error
	}
	res.Status = http.StatusOK
	res.Message = "Success"

	return res, nil
}

func NewPassword(user_id string, hashPassword string) (Response, error) {
	var res Response
	var user User

	db := config.GetDBInstance()
	result := db.Where("id = ?", user_id).Take(&user)
	if result.Error != nil {
		if is_notfound := errors.Is(result.Error, gorm.ErrRecordNotFound); is_notfound {
			res.Status = http.StatusOK
			res.Message = "can't find record"
			return res, result.Error
		}

		res.Status = http.StatusInternalServerError
		res.Message = "something went wrong"
		return res, result.Error
	}

	user.Password = hashPassword

	db.Save(&user)

	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = user

	return res, nil
}
