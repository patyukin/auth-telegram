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

func (uc *UseCase) SignUp(ctx context.Context, in model.SignUpData) (dto.SignUpResponse, error) {
	var err error
	in.Password, err = uc.HashPassword(in.Password)
	if err != nil {
		return dto.SignUpResponse{}, fmt.Errorf("failed to hash password: %w", err)
	}

	var userUUID, code uuid.UUID

	err = uc.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		userUUID, err = uc.registry.GetRepo().InsertIntoUser(ctx, in)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		_, err = uc.registry.GetRepo().InsertIntoTelegramUsers(ctx, in.Telegram, userUUID)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		code, err = uuid.NewV7FromReader(strings.NewReader(userUUID.String()))
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		err = uc.redis.SetSignUpCode(ctx, in.Telegram, code, userUUID, time.Hour)
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
