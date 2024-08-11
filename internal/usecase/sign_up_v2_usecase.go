package usecase

import (
	"auth-telegram/internal/db"
	"auth-telegram/internal/handler/dto"
	"auth-telegram/internal/model"
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (uc *UseCase) SignUpV2(ctx context.Context, in model.SignUpV2Data) (dto.SignUpV2Response, error) {
	var err error
	var userUUID uuid.UUID

	err = uc.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		in.Password, err = uc.HashPassword(in.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		userUUID, err = uc.registry.GetRepo().InsertIntoUserV2(ctx, in)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		return nil
	})
	if err != nil {
		return dto.SignUpV2Response{}, fmt.Errorf("failed to create user: %w", err)
	}

	return dto.SignUpV2Response{UserUUID: userUUID.String()}, nil
}
