package handler

import (
	"auth-telegram/internal/handler/dto"
	"auth-telegram/internal/model"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
)

const (
	HeaderAuthorization = "Authorization"
	HeaderUserID        = "User-ID"
	HeaderUserRole      = "User-Role"
)

type UseCase interface {
	SignUp(ctx context.Context, loginData model.SignUpData) (dto.SignUpResponse, error)
	ResendCode(ctx context.Context, loginData model.SignUpData) (dto.SignUpResponse, error)
	SignIn(ctx context.Context, signInData model.SignInData) error
	SignInVerify(ctx context.Context, code string) (*dto.TokensResponse, error)
	GetUserIDByToken(ctx context.Context, token string) (uuid.UUID, error)
	GetUserFullInfo(ctx context.Context, userID uuid.UUID) (model.User, error)
	GetTelegramBot() string
	GetJWTToken() string
}

type Handler struct {
	uc UseCase
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func New(uc UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) HandleError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	log.Error().Msgf("Error: %s", message)
	err := json.NewEncoder(w).Encode(ErrorResponse{Error: message})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
