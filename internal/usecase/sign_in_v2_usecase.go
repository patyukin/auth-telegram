package usecase

import (
	"auth-telegram/internal/handler/dto"
	"auth-telegram/internal/model"
	"context"
	"fmt"
)

func (uc *UseCase) SignInV2(ctx context.Context, signInData model.SignInData) (dto.TokensResponse, error) {
	user, err := uc.registry.GetRepo().SelectUserByLogin(ctx, signInData.Login)
	if err != nil {
		return dto.TokensResponse{}, fmt.Errorf("failed to select user: %w", err)
	}

	err = uc.ComparePasswords([]byte(user.PasswordHash), signInData.Password)
	if err != nil {
		return dto.TokensResponse{}, fmt.Errorf("failed to compare passwords: %w", err)
	}

	token, err := uc.generateJWT(user.UUID.String())
	if err != nil {
		return dto.TokensResponse{}, fmt.Errorf("failed to generate jwt: %w", err)
	}

	refreshToken, err := uc.registry.GetRepo().InsertToken(ctx, user.UUID)
	if err != nil {
		return dto.TokensResponse{}, fmt.Errorf("failed to insert token: %w", err)
	}

	return dto.TokensResponse{AccessToken: token, RefreshToken: refreshToken.String()}, nil
}
