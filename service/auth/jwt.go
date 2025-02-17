package auth

import (
	"context"
	"ecom/config"
	"ecom/types"
	"ecom/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type contextKey string
const UserKey contextKey = "userId"

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

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {


	return func(w http.ResponseWriter, r *http.Request) {
		//get token from user request
		tokenString := getTokenFromRequest(r)
		//validate JWT
		token, err := validateToken(tokenString)

		if err != nil {
			log.Printf("failed to vlidate token %v", err)
			permissionDenied(w)
			return
		}
		//fetch userId from database

		if !token.Valid {
			log.Printf("invlaid token")
			permissionDenied(w)
			return
		}
		//set userId in context

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userId, err := strconv.Atoi(str)

		u, err := store.GetUserById(userId)

		if err != nil {
			log.Printf("failed to get user by id %d", userId)
			permissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)

		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}


func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")

	if tokenAuth != "" {
		return tokenAuth
	}

	return ""
}

func validateToken(t string)(*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token)(interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
}

func GetUserIDFromContext(ctx context.Context) int {
	userId, ok := ctx.Value(UserKey).(int)

	if !ok {
		return -1
	}

	return userId
}