package converter

import (
	"auth-telegram/internal/handler/dto"
	"auth-telegram/internal/model"
	"time"
)

func ToUserFromModelUser(modelUser model.User) dto.User {
	result := dto.User{
		ID:        modelUser.UUID.String(),
		Login:     modelUser.Login,
		CreatedAt: modelUser.CreatedAt.String(),
	}

	result.Name = &modelUser.Name.String
	if !modelUser.Name.Valid {
		result.Name = nil
	}

	result.Surname = &modelUser.Surname.String
	if !modelUser.Surname.Valid {
		result.Surname = nil
	}

	updatedAt := modelUser.UpdatedAt.Time.Format(time.DateTime)
	result.UpdatedAt = &updatedAt
	if !modelUser.UpdatedAt.Valid {
		result.UpdatedAt = nil
	}

	return result
}

func ToUserAuthInfoFromModelUserInfo(modelUser model.UserAuthInfo) dto.UserAuthInfo {
	return dto.UserAuthInfo{
		UserUUID: modelUser.UserUUID.String(),
		Role:     string(modelUser.Role),
	}
}
