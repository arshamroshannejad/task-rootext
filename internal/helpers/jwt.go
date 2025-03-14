package helpers

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

func IsTokenValid(reqToken string, secretKey string) (*jwt.Token, error) {
	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func GetClaims(token *jwt.Token) (M, error) {
	claimsMap := make(M)
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		for k, v := range claims {
			claimsMap[k] = v
		}
		return claimsMap, nil
	}
	return nil, errors.New("invalid token claims")
}
