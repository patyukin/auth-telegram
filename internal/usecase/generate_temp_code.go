package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func (uc *UseCase) GenerateSignInCode() (string, error) {
	bytes := make([]byte, 30)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}
