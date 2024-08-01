package usecase

import (
	"auth-telegram/internal/model"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"time"
)

func (uc *UseCase) SignIn(ctx context.Context, signInData model.SignInData) error {
	user, err := uc.registry.GetRepo().SelectUserByLogin(ctx, signInData.Login)
	if err != nil {
		return fmt.Errorf("failed to select user: %w", err)
	}

	err = uc.ComparePasswords([]byte(user.PasswordHash), signInData.Password)
	if err != nil {
		return fmt.Errorf("failed to compare passwords: %w", err)
	}

	// Генерация уникального кода 2FA
	var code string
	var exists int64
	for {
		code, err = uc.GenerateSignInCode()
		if err != nil {
			return fmt.Errorf("failed to generate sign in code: %w", err)
		}

		// Проверка на уникальность кода
		exists, err = uc.redis.Exists2FACode(ctx, code)
		if err != nil {
			return fmt.Errorf("failed to check sign in code: %w", err)
		}

		if exists == 0 {
			break
		}
	}

	err = uc.redis.Set2FACode(ctx, code, user.UUID.String(), 5*time.Minute)
	if err != nil {
		return fmt.Errorf("failed to set 2fa code: %w", err)
	}

	telegramUser, err := uc.registry.GetRepo().SelectFromTelegramUsersByUser(ctx, user.UUID)
	if err != nil {
		return fmt.Errorf("failed to select user: %w", err)
	}

	msg := tgbotapi.NewMessage(telegramUser.TgChatID, code)
	if _, err = uc.bot.API.Send(msg); err != nil {
		log.Error().Msgf("Error sending message: %v", err)
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Info().Msgf("Sent message to chat %d", telegramUser.TgChatID)
	return nil
}
