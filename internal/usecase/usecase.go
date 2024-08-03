package usecase

import (
	"auth-telegram/internal/cacher"
	"auth-telegram/internal/converter"
	"auth-telegram/internal/db"
	"auth-telegram/internal/handler/dto"
	"auth-telegram/internal/telegram"
	"context"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	registry  *db.Client
	redis     *cacher.Cacher
	bot       *telegram.Bot
	jwtSecret []byte
}

func New(registry *db.Client, redis *cacher.Cacher, bot *telegram.Bot, jwtSecret string) *UseCase {
	return &UseCase{
		registry:  registry,
		redis:     redis,
		bot:       bot,
		jwtSecret: []byte(jwtSecret),
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

func (uc *UseCase) GetUserAuthInfoByToken(ctx context.Context, id string) (dto.UserAuthInfo, error) {
	userAuthInfo, err := uc.registry.GetRepo().SelectUserAuthInfoByUUID(ctx, id)
	if err != nil {
		return dto.UserAuthInfo{}, err
	}

	return converter.ToUserAuthInfoFromModelUserInfo(userAuthInfo), nil
}

func (uc *UseCase) GetJWTToken() []byte {
	return uc.jwtSecret
}

func (uc *UseCase) GetTelegramBot() string {
	return uc.bot.API.Self.UserName
}
