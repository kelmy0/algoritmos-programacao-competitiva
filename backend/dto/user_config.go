package dto

import "time"

type ChangePasswordRequest struct {
	OldPassword        string `json:"old_password" binding:"required,min=8"`
	NewPassword        string `json:"new_password" binding:"required,min=8"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required,min=8"`
}

type ChangePasswordResponse struct {
	Code                   string `json:"code"`
	Message                string `json:"message"`
	OthersDevicesLoggedOut bool   `json:"others_devices_logged_out"`
}

type DefinePasswordRequest struct {
	NewPassword        string `json:"new_password" binding:"required,min=8"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required,min=8"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,min=5"`
}

type ResetPasswordRequest struct {
	Token              string `json:"token" binding:"required"`
	NewPassword        string `json:"new_password" binding:"required,min=8"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required,min=8"`
}

type GetMyCredentialsResponse struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	RoleId    int        `json:"role_id"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
