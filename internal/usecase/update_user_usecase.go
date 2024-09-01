package usecase

import (
	"auth-telegram/internal/handler/dto"
	"context"
	"fmt"
)

func (uc *UseCase) UpdateUser(ctx context.Context, user dto.UpdateUserRequest, userUUID string) error {
	err := uc.registry.GetRepo().UpdateUser(ctx, user, userUUID)
	if err != nil {
		return fmt.Errorf("failed updating user: %w", err)
	}

	return nil
}
