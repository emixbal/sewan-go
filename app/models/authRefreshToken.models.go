package models

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/golang-jwt/jwt"
)

type UserClaim struct {
	Issuer  string
	Id      int
	Email   string
	IsAdmin bool
}

var (
	jwtKey        = []byte(os.Getenv("JWT_SECRET"))
	jwtRefreshKey = []byte(os.Getenv("REFRESH_SECRET"))
)

// GenerateTokens returns the access and refresh tokens
func GenerateTokens(userClaim *UserClaim, isDoRefresh bool) (string, string) {
	issuer := userClaim.Issuer
	if isDoRefresh {
		issuer = utils.UUIDv4()
	}

	accessToken := GenerateAccessClaims(userClaim, issuer)
	refreshToken := GenerateRefreshClaims(userClaim, issuer)

	return accessToken, refreshToken
}

// GenerateAccessClaims returns a claim and a acess_token string
func GenerateAccessClaims(userClaim *UserClaim, issuer string) string {
	claim := &jwt.MapClaims{
		"issuer":   issuer,
		"email":    userClaim.Email,
		"user_id":  userClaim.Id,
		"is_admin": userClaim.IsAdmin,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	return tokenString
}

// GenerateRefreshClaims returns refresh_token
func GenerateRefreshClaims(userClaim *UserClaim, issuer string) string {
	refreshClaim := &jwt.MapClaims{
		"issuer":   issuer,
		"email":    userClaim.Email,
		"user_id":  userClaim.Id,
		"is_admin": userClaim.IsAdmin,
		"exp":      time.Now().Add(30 * 24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)
	refreshTokenString, err := refreshToken.SignedString(jwtRefreshKey)
	if err != nil {
		panic(err)
	}

	return refreshTokenString
}
