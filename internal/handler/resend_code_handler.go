package handler

import (
	"auth-telegram/internal/handler/dto"
	"auth-telegram/internal/model"
	"auth-telegram/pkg/httperror"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

// ResendCodeHandler godoc
// @Summary      Повторная отправка кода
// @Description  Повторная отправка кода для окончания регистрации
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body model.SignUpData true "Sign-up data required to resend the code"
// @Success      200   {object}  dto.ResendCodeResponse  "Message with link to complete registration"
// @Failure      400   {object}  ErrorResponse    "Invalid input data"
// @Failure      500   {object}  ErrorResponse    "Internal server error"
// @Router       /resend-code [post]
func (h *Handler) ResendCodeHandler(w http.ResponseWriter, r *http.Request) {
	var in model.SignUpData

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httperror.SendError(w, "invalid data", http.StatusBadRequest)
		return
	}

	if !in.IsValid() {
		log.Error().Msg("invalid input")
		httperror.SendError(w, "invalid input", http.StatusBadRequest)
		return
	}

	signUpResponse, err := h.uc.ResendCode(r.Context(), in)
	if err != nil {
		log.Error().Msgf("invalid input, err: %v", err)
		httperror.SendError(w, "invalid sign up", http.StatusBadRequest)
		return
	}

	link := fmt.Sprintf("/start=%s", signUpResponse.Code)
	message := fmt.Sprintf("Please send this message to complete your registration: %s", link)

	response := dto.ResendCodeResponse{Message: message}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		httperror.SendError(w, err.Error(), http.StatusInternalServerError)
	}
}
