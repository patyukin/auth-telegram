package handler

import (
	"auth-telegram/internal/handler/dto"
	"auth-telegram/internal/metrics"
	"auth-telegram/internal/model"
	"auth-telegram/pkg/httperror"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

// SignUpHandler godoc
// @Summary      Регистрация нового пользователя
// @Description  Регистрация нового пользователя
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body model.SignUpData true "Sign Up Data"
// @Success      200   {object}  dto.SignUpV1Response "1 этап регистрации завершен успешно"
// @Failure      400   {object}  ErrorResponse "Invalid request body"
// @Failure      500   {object}  ErrorResponse "Internal server error"
// @Router       /sign-up [post]
// @Router       /v1/sign-up [post]
func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
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

	signUp, err := h.uc.SignUp(r.Context(), in)
	if err != nil {
		log.Error().Msgf("invalid input, err: %v", err)
		httperror.SendError(w, "invalid sign up", http.StatusBadRequest)
		return
	}

	link := fmt.Sprintf("https://t.me/%s?start=%s", signUp.BotName, signUp.Code)
	response := dto.SignUpV1Response{
		Message: fmt.Sprintf("1 час для окончания регистрации. Пожалуйста, перейдите по ссылке в telegram бот и нажмите /start для окончания регистрации: %s", link),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	metrics.SignUpV1RegisterTraffic.Inc()
}
