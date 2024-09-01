package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

// GetUserInfoHandler godoc
// @Summary      Get user information
// @Description  Получение информации о пользователе по его UUID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user_id  header  string  true  "User ID (UUID)"
// @Success      200  {object}  dto.User
// @Failure      400  {object}  ErrorResponse  "Bad Request"
// @Failure      500  {object}  ErrorResponse  "Internal Server Error"
// @Router       /user/info [get]
func (h *Handler) GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.Header.Get(HeaderUserID))
	if err != nil {
		h.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	userUUID, err := uuid.Parse(id.String())
	if err != nil {
		h.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := h.uc.GetUserInfoByUUID(r.Context(), userUUID)
	if err != nil {
		h.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userInfo)
	if err != nil {
		h.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
