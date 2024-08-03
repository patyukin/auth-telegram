package usecase

import (
	"auth-telegram/internal/converter"
	"auth-telegram/internal/handler/dto"
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (uc *UseCase) GetUserInfoByUUID(ctx context.Context, userID uuid.UUID) (dto.User, error) {
	user, err := uc.registry.GetRepo().SelectUserByUUID(ctx, userID.String())
	if err != nil {
		return dto.User{}, fmt.Errorf("failed getting user: %w", err)
	}

	return converter.ToUserFromModelUser(user), nil
}
