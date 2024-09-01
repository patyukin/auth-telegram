package handler

import (
	"encoding/json"
	"net/http"
)

// AdminGetUserInfoHandler godoc
// @Summary      Get admin user information
// @Description  Получение информации о пользователе по его UUID для админа
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user_id  header  string  true  "User ID (UUID)"
// @Success      200  {object}  dto.AdminUserInfo
// @Failure      400  {object}  ErrorResponse  "Bad Request"
// @Failure      500  {object}  ErrorResponse  "Internal Server Error"
// @Router       /admin/user-info/{id} [get]
func (h *Handler) AdminGetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get(HeaderUserID)
	if id == "" {
		h.HandleError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	role := r.Header.Get(HeaderUserRole)
	if role != "admin" {
		h.HandleError(w, http.StatusForbidden, "forbidden")
		return
	}

	userInfo, err := h.uc.GetUserAuthInfoByAdmin(r.Context(), id)
	if err != nil {
		h.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(userInfo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
