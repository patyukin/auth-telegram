package handler

import (
	"auth-telegram/internal/metrics"
	"auth-telegram/internal/model"
	"auth-telegram/pkg/httperror"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

// SignUpV2Handler godoc
// @Summary      Регистрация нового пользователя
// @Description  Регистрация нового пользователя с версией v2
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body model.SignUpV2Data true "Sign Up v2 Data Request"
// @Success      200   {object}  dto.SignUpV2Response "Validate token successfully"
// @Failure      400   {object}  ErrorResponse "Invalid request body"
// @Failure      500   {object}  ErrorResponse "Internal server error"
// @Router       /v2/sign-up [post]
func (h *Handler) SignUpV2Handler(w http.ResponseWriter, r *http.Request) {
	var in model.SignUpV2Data
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := r.Header.Get("Authorization")
	tokenFromDB, err := h.uc.GetTokenByName(r.Context(), "recipe")
	if err != nil {
		log.Error().Err(err).Msgf("failed to sign up, error: %v", err)
		httperror.SendError(w, "invalid sign up", http.StatusBadRequest)
		return
	}

	if token != tokenFromDB {
		log.Error().Err(err).Msgf("failed to sign up, error: %v", err)
		httperror.SendError(w, "invalid sign up", http.StatusBadRequest)
		return
	}

	result, err := h.uc.SignUpV2(r.Context(), in)
	if err != nil {
		log.Error().Err(err).Msgf("failed to sign up, error: %v", err)
		httperror.SendError(w, "invalid sign up", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(result)
	if err != nil {
		log.Error().Err(err).Msgf("failed to sign up, error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(res)
	if err != nil {
		log.Error().Err(err).Msgf("failed to sign up, error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metrics.SignUpV2RegisterTraffic.Inc()
}
