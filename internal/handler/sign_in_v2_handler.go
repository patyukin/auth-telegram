package handler

import (
	"auth-telegram/internal/metrics"
	"auth-telegram/internal/model"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"net/http"
)

// SignInV2Handler godoc
// @Summary      Вход в систему с версией v2
// @Description  Аутентификация пользователя по логину и паролю в системе с версией v2
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body  model.SignInData  true  "Sign-in data"
// @Success      200   {object}  dto.TokensResponse  "Successfully signed in"
// @Failure      400   {object}  ErrorResponse  "Invalid input data"
// @Failure      500   {object}  ErrorResponse  "Failed to sign in"
// @Router       /v2/sign-in [post]
func (h *Handler) SignInV2Handler(w http.ResponseWriter, r *http.Request) {
	var err error

	_, span := otel.Tracer("auth-telegram").Start(r.Context(), "SignInV2")
	timer := prometheus.NewTimer(metrics.SignInLatencyHistogram)
	metrics.TotalAuthentications.Inc()

	defer func() {
		span.End()
		timer.ObserveDuration()
		if err != nil {
			metrics.FailedAuthentications.Inc()
		} else {
			metrics.SuccessfulAuthentications.Inc()
		}
	}()

	var signInData model.SignInData
	if err = json.NewDecoder(r.Body).Decode(&signInData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokens, err := h.uc.SignInV2(r.Context(), signInData)
	if err != nil {
		log.Error().Err(err).Msgf("failed to sign in, error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(tokens); err != nil {
		log.Error().Err(err).Msgf("failed to encode tokens, error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
