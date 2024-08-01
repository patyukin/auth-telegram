package dto

type SignUpResponse struct {
	BotName string `json:"bot_name"`
	Code    string `json:"code"`
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
