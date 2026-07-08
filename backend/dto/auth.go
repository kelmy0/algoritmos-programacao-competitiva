package dto

type AuthRequest struct {
	Email    string `json:"email" binding:"required,min=5"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	Requires2FA  bool   `json:"requires_2fa"`
	PreAuthToken string `json:"pre_auth_token,omitempty"`
}

type Verify2FARequest struct {
	PreAuthToken string `json:"pre_auth_token" binding:"required"`
	Code         string `json:"code" binding:"required,len=6"`
}
