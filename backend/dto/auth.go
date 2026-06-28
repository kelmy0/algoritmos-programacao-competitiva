package dto

type AuthRequest struct {
	Email    string `json:"email" binding:"required,min=5"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	AcessToken   string `json:"acess_token"`
	RefreshToken string `json:"refresh_token"`
}
