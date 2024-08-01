package usecase

import (
	"context"
	"fmt"
)

func (uc *UseCase) CleanTokens(ctx context.Context) error {
	err := uc.registry.GetRepo().CleanTokens(ctx)
	if err != nil {
		return fmt.Errorf("failed cleaning tokens: %w", err)
	}

	return nil
}
