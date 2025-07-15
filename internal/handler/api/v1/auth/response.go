package auth

type tokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type userIDResponse struct {
	UserID string `json:"user_id"`
}
