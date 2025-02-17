package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"houseparty.com/config"
)

func GenerateToken(email, username string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email": email,
		"userId": userId,
        "username": username,
        "exp": time.Now().Add(time.Hour*2).Unix(),
    })

    return token.SignedString([]byte(config.GetJWTKey()))
}

func VerifyToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error){
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(config.GetJWTKey()), nil
	})
	if err != nil {
		return 0, errors.New("parsing jwt error")
	}

	tokenIsValid := parsedToken.Valid
	if !tokenIsValid {
		return 0, errors.New("token invalid")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	userId := claims["userId"].(float64)
	n := int64(userId) 
	
	return n, nil
}