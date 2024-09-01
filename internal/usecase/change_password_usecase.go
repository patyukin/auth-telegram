package usecase

import (
	"auth-telegram/internal/db"
	"auth-telegram/internal/handler/dto"
	"auth-telegram/internal/model"
	"context"
	"fmt"
)

func (uc *UseCase) ChangePassword(ctx context.Context, req dto.ChangePasswordRequest, userUUID string) error {
	var err error
	var passwordData model.CheckerPasswordData

	err = uc.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		passwordData, err = uc.registry.GetRepo().SelectUserPasswordDataByUUID(ctx, userUUID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		err = uc.ComparePasswords([]byte(passwordData.PasswordHash), req.OldPassword)
		if err != nil {
			return fmt.Errorf("failed to compare passwords: %w", err)
		}

		req.NewPassword, err = uc.HashPassword(req.NewPassword)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		err = uc.registry.GetRepo().UpdateUserPassword(ctx, req, userUUID)
		if err != nil {
			return fmt.Errorf("failed to change password: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}

	return nil
}
