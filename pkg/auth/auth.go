package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var secretKey = os.Getenv("SECRET_KEY_TOKEN")

func GetToken(userName string, password string, isAdmin bool) (token string, err error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name": userName,
		"password":  password,
		"is_admin":  isAdmin,
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})
	if token, err = tokenClaims.SignedString([]byte(secretKey)); err != nil {
		return
	}
	return
}

func ParseToken(tokenStr string) (claims jwt.MapClaims, err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return
	}
	if c, ok := token.Claims.(jwt.MapClaims); ok {
		claims = c
		return
	}
	err = fmt.Errorf("ParseToken error")
	return
}
