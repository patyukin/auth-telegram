package dto

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ResendCodeResponse struct {
	Message string `json:"message"`
}

type SignUpV1Response struct {
	Message string `json:"message"`
}

type SignUpResponse struct {
	BotName string `json:"bot_name"`
	Code    string `json:"code"`
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	ID        string  `json:"id"`
	Login     string  `json:"login"`
	Name      *string `json:"name"`
	Surname   *string `json:"surname"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

type UpdateUserRequest struct {
	Name    *string `json:"name"`
	Surname *string `json:"surname"`
}

type UserAuthInfo struct {
	UserUUID string `json:"user_id"`
	Role     string `json:"role"`
}

type SignUpV2Response struct {
	UserUUID string `json:"user_id"`
}

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	UUID string `json:"id"`
}

type AdminUserInfo struct {
	ID       string  `json:"id"`
	Login    string  `json:"login"`
	Name     *string `json:"name"`
	Surname  *string `json:"surname"`
	Role     string  `json:"role"`
	Telegram struct {
		Username string `json:"username"`
	} `json:"telegram"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
