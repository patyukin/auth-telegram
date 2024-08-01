package usecase

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func (uc *UseCase) ValidateTempCodeAndGenerateToken(ctx context.Context, userID, tempCode, username, password string) (string, error) {
	storedCode, err := uc.redis.GetTempCode(ctx, userID)
	if err != nil || storedCode != tempCode {
		return "", err
	}

	err = uc.redis.DeleteTempCode(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to delete temp code: %w", err)
	}

	if username != "expectedUsername" || password != "expectedPassword" {
		return "", fmt.Errorf("invalid username or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(1 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
