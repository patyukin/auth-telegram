package usecase

import (
	"auth-telegram/internal/model"
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (uc *UseCase) GetUserFullInfo(ctx context.Context, userID uuid.UUID) (model.User, error) {
	user, err := uc.registry.GetRepo().SelectUserByUUID(ctx, userID.String())
	if err != nil {
		return model.User{}, fmt.Errorf("failed getting user: %w", err)
	}

	return user, nil
}
