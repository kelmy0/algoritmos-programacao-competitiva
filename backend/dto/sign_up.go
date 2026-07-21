package dto

type SignUpRequest struct {
	Name            string `json:"name" binding:"required,min=6,max=128"`
	Username        string `json:"username" binding:"required,min=6,max=32"`
	Email           string `json:"email" binding:"required,min=5,max=128"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

type SignUpResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	Success     bool   `json:"success"`
	AutoLogin   bool   `json:"auto_login"`
}
