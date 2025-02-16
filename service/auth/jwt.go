package auth

import (
	"ecom/config"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"

)

func CreateJWT(secret []byte, userId int)(string, error) {

	expiration := time.Second * time.Duration(config.Envs.JWTWxpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID" : strconv.Itoa(userId),
		"expiresAt" : time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}