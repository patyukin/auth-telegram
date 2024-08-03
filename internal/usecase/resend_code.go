package usecase

import (
	"auth-telegram/internal/db"
	"auth-telegram/internal/handler/dto"
	"auth-telegram/internal/model"
	"context"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

func (uc *UseCase) ResendCode(ctx context.Context, in model.SignUpData) (dto.SignUpResponse, error) {
	var err error
	var user model.CheckerPasswordData
	var code uuid.UUID

	err = uc.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		user, err = uc.registry.GetRepo().SelectUserByLogin(ctx, in.Login)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		code, err = uuid.NewV7FromReader(strings.NewReader(user.UUID.String()))
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		err = uc.redis.SetSignUpCode(ctx, in.Telegram, code, user.UUID, time.Minute*10)
		if err != nil {
			return fmt.Errorf("failed to register user: %w", err)
		}

		return nil
	})
	if err != nil {
		return dto.SignUpResponse{}, fmt.Errorf("failed to create user: %w", err)
	}

	return dto.SignUpResponse{BotName: uc.GetTelegramBot(), Code: code.String()}, nil
}
