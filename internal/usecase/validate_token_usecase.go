package usecase

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

func (uc *UseCase) ValidateToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return uc.GetJWTToken(), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("token is invalid")
	}

	id := token.Claims.(jwt.MapClaims)["id"].(string)

	return id, nil
}
