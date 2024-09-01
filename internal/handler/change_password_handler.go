package handler

import (
	"auth-telegram/internal/handler/dto"
	"encoding/json"
	"net/http"
)

// ChangePasswordHandler godoc
// @Summary      Обновление пароля
// @Description  Обновить пароль пользователя
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body dto.ChangePasswordRequest true "Passwords Request"
// @Success      200   {nil}  nil "Password updated successfully"
// @Failure      400   {object}  ErrorResponse "Invalid request body"
// @Failure      500   {object}  ErrorResponse "Internal server error"
// @Router       /refresh [post]
func (h *Handler) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	userUUID := r.Header.Get(HeaderUserID)
	if userUUID == "" {
		h.HandleError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := h.uc.ChangePassword(r.Context(), req, userUUID)
	if err != nil {
		h.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
