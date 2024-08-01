package usecase

import (
	"auth-telegram/internal/cacher"
	"auth-telegram/internal/db"
	"auth-telegram/internal/telegram"
	"context"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	registry  *db.Client
	redis     *cacher.Cacher
	bot       *telegram.Bot
	jwtSecret string
}

func New(registry *db.Client, redis *cacher.Cacher, bot *telegram.Bot, jwtSecret string) *UseCase {
	return &UseCase{
		registry:  registry,
		redis:     redis,
		bot:       bot,
		jwtSecret: jwtSecret,
	}
}

func (uc *UseCase) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (uc *UseCase) ComparePasswords(hashedPassword []byte, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(plainPassword))
	return err
}

func (uc *UseCase) CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (uc *UseCase) GetUserIDByToken(ctx context.Context, id string) (uuid.UUID, error) {
	userUUID, err := uc.registry.GetRepo().SelectUserUUIDByUUID(ctx, id)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userUUID, nil
}

func (uc *UseCase) GetJWTToken() string {
	return uc.jwtSecret
}

func (uc *UseCase) GetTelegramBot() string {
	return uc.bot.API.Self.UserName
}
