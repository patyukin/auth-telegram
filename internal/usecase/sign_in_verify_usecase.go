package usecase

import (
	"auth-telegram/internal/handler/dto"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func (uc *UseCase) SignInVerify(ctx context.Context, code string) (dto.TokensResponse, error) {
	userID, err := uc.redis.Get2FACode(ctx, code)
	if err != nil {
		return dto.TokensResponse{}, fmt.Errorf("failed to get 2fa code: %w", err)
	}

	err = uc.redis.Delete2FACode(ctx, code)
	if err != nil {
		return dto.TokensResponse{}, fmt.Errorf("failed to delete 2fa code: %w", err)
	}

	user, err := uc.registry.GetRepo().SelectUserByUUID(ctx, userID)
	if err != nil {
		return dto.TokensResponse{}, fmt.Errorf("failed to get user: %w", err)
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

func (uc *UseCase) generateJWT(userUUID string) (string, error) {
	claims := jwt.MapClaims{
		"id":  userUUID,
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(uc.jwtSecret))
}
