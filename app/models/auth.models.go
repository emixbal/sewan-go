package models

import (
	"errors"
	"fmt"
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

func CheckLogin(email, passwordTxt string) (isExist bool, isMatch bool, tokenString string, userObj User, errMessage error) {
	var user User
	db := config.GetDBInstance()

	if result := db.Where(&User{Email: email}).First(&user); result.Error != nil {
		if is_notfound := errors.Is(result.Error, gorm.ErrRecordNotFound); is_notfound {
			return false, false, "", user, result.Error
		}

		fmt.Print(result.Error)
		return true, false, "", user, result.Error
	}

	match, _ := helpers.CheckPasswordHash(user.Password, passwordTxt)
	if !match {
		return true, false, "", user, nil
	}

	// atClaims := jwt.MapClaims{}
	// atClaims["authorized"] = true
	// atClaims["user_id"] = user.ID
	// atClaims["email"] = user.Email
	// atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	// at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	// tokenString, err := at.SignedString([]byte(os.Getenv("JWT_SECRET")))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    email,
		"user_id":  user.ID,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
		"iat":      time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		fmt.Println(err)
	}

	return true, true, tokenString, user, nil
}
