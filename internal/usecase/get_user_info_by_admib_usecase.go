package usecase

import (
	"auth-telegram/internal/handler/dto"
	"context"
	"fmt"
	"time"
)

func (uc *UseCase) GetUserAuthInfoByAdmin(ctx context.Context, id string) (dto.AdminUserInfo, error) {
	userInfo, err := uc.registry.GetRepo().SelectAdminUserByUserUUID(ctx, id)
	if err != nil {
		return dto.AdminUserInfo{}, fmt.Errorf("failed to get user info: %w", err)
	}

	userInfoDto := dto.AdminUserInfo{
		ID:        userInfo.UUID.String(),
		Login:     userInfo.Login,
		Role:      string(userInfo.Role),
		CreatedAt: userInfo.CreatedAt.Format(time.DateTime),
	}

	if userInfo.UpdatedAt.Valid {
		updatedAt := userInfo.UpdatedAt.Time.Format(time.DateTime)
		userInfoDto.UpdatedAt = &updatedAt
	}

	userInfoDto.Name = &userInfo.Name.String
	userInfoDto.Surname = &userInfo.Surname.String
	userInfoDto.Telegram.Username = userInfo.Telegram.Username

	return userInfoDto, nil
}
