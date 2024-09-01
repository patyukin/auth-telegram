package handler

import (
	"auth-telegram/internal/model"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

// SignInHandler godoc
// @Summary      Вход в систему
// @Description  Аутентификация пользователя по логину и паролю в системе
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body  model.SignInData  true  "Sign-in data"
// @Success      200   {nil}  nil  "Successfully signed in, no content"
// @Failure      400   {object}  ErrorResponse  "Invalid input data"
// @Failure      500   {object}  ErrorResponse  "Failed to sign in"
// @Router       /sign-in [post]
func (h *Handler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	var signInData model.SignInData

	if err := json.NewDecoder(r.Body).Decode(&signInData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.uc.SignIn(r.Context(), signInData)
	if err != nil {
		log.Error().Err(err).Msgf("failed to sign in, error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
