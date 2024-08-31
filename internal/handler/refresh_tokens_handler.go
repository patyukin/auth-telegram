package handler

import (
	"auth-telegram/internal/handler/dto"
	"encoding/json"
	"net/http"
)

// GenerateRefreshTokensHandler godoc
// @Summary      Generate new access and refresh tokens
// @Description  Generates new access and refresh tokens using a provided refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body dto.RefreshTokenRequest true "Refresh Token Request"
// @Success      200   {object}  dto.TokensResponse "Tokens generated successfully"
// @Failure      400   {object}  ErrorResponse "Invalid request body"
// @Failure      500   {object}  ErrorResponse "Internal server error"
// @Router       /auth/refresh-tokens [post]
func (h *Handler) GenerateRefreshTokensHandler(w http.ResponseWriter, r *http.Request) {
	var refreshToken dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshToken); err != nil {
		h.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.uc.GenerateTokens(r.Context(), refreshToken.RefreshToken)
	if err != nil {
		h.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err = json.NewEncoder(w).Encode(tokens); err != nil {
		h.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
