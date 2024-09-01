package handler

import (
	"auth-telegram/internal/handler/dto"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

// SetUserInfoHandler godoc
// @Summary      Обновление информации о пользователе
// @Description  Пользователь сам обновляет информацию о себе
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body  dto.UpdateUserRequest  true  "Sign-in data"
// @Success      200   {nil}  nil  "Successfully updated info"
// @Failure      400   {object}  ErrorResponse  "Invalid input data"
// @Failure      500   {object}  ErrorResponse  "Failed to update user info"
// @Router       /sign-in [post]
func (h *Handler) SetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	var userDto dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userDto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userUUID := r.Header.Get("user_id")
	if userUUID == "" {
		log.Error().Msgf("user_id is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.uc.UpdateUser(r.Context(), userDto, userUUID); err != nil {
		log.Error().Err(err).Msgf("failed to set user info, error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
