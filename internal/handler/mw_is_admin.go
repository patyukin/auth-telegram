package handler

import (
	"auth-telegram/internal/model"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := r.Header.Get(HeaderUserID)
		userID, err := uuid.Parse(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			var userInfo model.User
			userInfo, err = h.uc.GetUserFullInfo(r.Context(), userID)
			if err != nil {
				return
			}

			if userInfo.Role != model.UserRoleAdmin {
				return
			}

			next.ServeHTTP(w, r)
		}
	})
}
