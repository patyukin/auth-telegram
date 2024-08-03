package usecase

import (
	"auth-telegram/internal/db"
	"auth-telegram/internal/handler/dto"
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (uc *UseCase) GenerateTokens(ctx context.Context, refreshToken string) (dto.TokensResponse, error) {
	var err error
	var token uuid.UUID
	var accessToken string
	var userUUID uuid.UUID

	err = uc.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		userUUID, err = uc.registry.GetRepo().GetUserUUIDByRefreshToken(ctx, refreshToken)
		if err != nil {
			return fmt.Errorf("failed generating tokens: %w", err)
		}

		err = uc.registry.GetRepo().DeleteToken(ctx, refreshToken)
		if err != nil {
			return fmt.Errorf("failed generating tokens: %w", err)
		}

		token, err = uc.registry.GetRepo().InsertToken(ctx, userUUID)
		if err != nil {
			return fmt.Errorf("failed generating tokens: %w", err)
		}

		accessToken, err = uc.generateJWT(userUUID.String())
		if err != nil {
			return fmt.Errorf("failed generating tokens: %w", err)
		}

		return nil
	})
	if err != nil {
		return dto.TokensResponse{}, fmt.Errorf("failed generating tokens: %w", err)
	}

	return dto.TokensResponse{AccessToken: accessToken, RefreshToken: token.String()}, nil
}
