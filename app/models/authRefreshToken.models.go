package models

import (
	"os"
	"sejuta-cita/config"
	"time"

	"github.com/golang-jwt/jwt"
)

// Claims represent the structure of the JWT token
type Claims struct {
	jwt.StandardClaims
	ID uint `gorm:"primaryKey"`
}

var jwtKey = []byte(os.Getenv("JWT_SECRET"))
var jwtRefreshKey = []byte(os.Getenv("REFRESH_SECRET"))

// GenerateTokens returns the access and refresh tokens
func GenerateTokens(email string) (string, string) {
	claim, accessToken := GenerateAccessClaims(email)
	refreshToken := GenerateRefreshClaims(claim)

	return accessToken, refreshToken
}

// GenerateAccessClaims returns a claim and a acess_token string
func GenerateAccessClaims(email string) (*Claims, string) {

	t := time.Now()
	claim := &Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    email,
			ExpiresAt: t.Add(15 * time.Minute).Unix(),
			Subject:   "access_token",
			IssuedAt:  t.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	return claim, tokenString
}

// GenerateRefreshClaims returns refresh_token
func GenerateRefreshClaims(cl *Claims) string {
	db := config.GetDBInstance()

	result := db.Where(&Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer: cl.Issuer,
		},
	}).Find(&Claims{})

	// checking the number of refresh tokens stored.
	// If the number is higher than 3, remove all the refresh tokens and leave only new one.
	if result.RowsAffected > 3 {
		db.Where(&Claims{
			StandardClaims: jwt.StandardClaims{Issuer: cl.Issuer},
		}).Delete(&Claims{})
	}

	t := time.Now()
	refreshClaim := &Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    cl.Issuer,
			ExpiresAt: t.Add(30 * 24 * time.Hour).Unix(),
			Subject:   "refresh_token",
			IssuedAt:  t.Unix(),
		},
	}

	// create a claim on DB
	db.Create(&refreshClaim)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)
	refreshTokenString, err := refreshToken.SignedString(jwtRefreshKey)
	if err != nil {
		panic(err)
	}

	return refreshTokenString
}
