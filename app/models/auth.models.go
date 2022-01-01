package models

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sejuta-cita/app/helpers"
	"sejuta-cita/config"
	"time"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func UserRegister(user *User) (Response, error) {
	var res Response
	db := config.GetDBInstance()

	result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)
	if result.Error != nil {
		fmt.Print("error CreateAUser")
		fmt.Print(result.Error)

		res.Status = http.StatusInternalServerError
		res.Message = "error save new record"
		return res, result.Error
	}

	res.Status = http.StatusOK
	res.Message = "success"

	return res, nil
}

func UserUpdate(user *User) (Response, error) {
	var res Response
	db := config.GetDBInstance()

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "email", "password"}),
	}).Create(&user)

	if result.Error != nil {
		log.Println("error CreateAUser")
		log.Println(result.Error)

		res.Status = http.StatusInternalServerError
		res.Message = "error save new record"
		return res, result.Error
	}

	res.Status = http.StatusOK
	res.Message = "success"

	return res, nil
}

func CheckLogin(email, passwordTxt string) (bool, bool, string, error) {
	var user User
	db := config.GetDBInstance()

	if result := db.Where(&User{Email: email}).First(&user); result.Error != nil {
		if is_notfound := errors.Is(result.Error, gorm.ErrRecordNotFound); is_notfound {
			return false, false, "", result.Error
		}

		fmt.Print(result.Error)
		return true, false, "", result.Error
	}

	match, _ := helpers.CheckPasswordHash(user.Password, passwordTxt)
	if !match {
		return true, false, "", nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    email,
		"userId":   user.ID,
		"is_admin": user.IsAdmin,
		"nbf":      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		fmt.Println(err)
	}

	return true, true, tokenString, nil
}
