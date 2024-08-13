package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
)

func (uc *UseCase) GetTokenByName(ctx context.Context, name string) (string, error) {
	token, err := uc.registry.GetRepo().SelectServiceTokenByName(ctx, name)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	h := sha256.New()
	_, _ = h.Write([]byte(token))
	token = fmt.Sprintf("%x", h.Sum(nil))

	return token, nil
}
