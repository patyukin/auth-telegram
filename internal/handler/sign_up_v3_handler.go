package handler

import (
	"auth-telegram/internal/model"
	"auth-telegram/pkg/httperror"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

// SignUpV3Handler godoc
// @Summary      Регистрация нового пользователя
// @Description  Регистрация нового пользователя с версией v3
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body model.SignUpV2Data true "Sign Up v3 Data Request"
// @Success      200   {object}  dto.SignUpV2Response "Validate token successfully"
// @Failure      400   {object}  ErrorResponse "Invalid request body"
// @Failure      500   {object}  ErrorResponse "Internal server error"
// @Router       /v3/sign-up [post]
func (h *Handler) SignUpV3Handler(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msgf("start sign up")

	var in model.SignUpV2Data
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenWithService := r.Header.Get("Authorization")
	parts := strings.SplitN(tokenWithService, ":", 2)

	if len(parts) != 2 {
		httperror.SendError(w, "invalid sign up", http.StatusBadRequest)
		return
	}

	service := parts[0]
	token := parts[1]

	tokenFromDB, err := h.uc.GetTokenByName(r.Context(), service)
	if err != nil {
		log.Error().Err(err).Msgf("failed to sign up in h.uc.GetTokenByName, error: %v", err)
		httperror.SendError(w, "invalid sign up", http.StatusBadRequest)
		return
	}

	if token != tokenFromDB {
		log.Error().Err(err).Msgf("failed to sign up tokens do not match, error: %v", err)
		httperror.SendError(w, "invalid sign up", http.StatusBadRequest)
		return
	}

	result, err := h.uc.SignUpV2(r.Context(), in)
	if err != nil {
		log.Error().Err(err).Msgf("failed to sign up in h.uc.SignUpV2, error: %v", err)
		httperror.SendError(w, "invalid sign up", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
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
}
