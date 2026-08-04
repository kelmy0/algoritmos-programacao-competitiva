package dto

import "time"

type ChangePasswordRequest struct {
	OldPassword        string `json:"oldPassword" binding:"required,min=8"`
	NewPassword        string `json:"newPassword" binding:"required,min=8"`
	ConfirmNewPassword string `json:"confirmNewPassword" binding:"required,min=8"`
}

type ChangePasswordResponse struct {
	Code                   string `json:"code"`
	Message                string `json:"message"`
	OthersDevicesLoggedOut bool   `json:"othersDevicesLoggedOut"`
}

type DefinePasswordRequest struct {
	NewPassword        string `json:"newPassword" binding:"required,min=8"`
	ConfirmNewPassword string `json:"confirmNewPassword" binding:"required,min=8"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,min=5"`
}

type ResetPasswordRequest struct {
	Token              string `json:"token" binding:"required"`
	NewPassword        string `json:"newPassword" binding:"required,min=8"`
	ConfirmNewPassword string `json:"confirmNewPassword" binding:"required,min=8"`
}

type GetMyCredentialsResponse struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	RoleId    int        `json:"roleId"`
	LastLogin *time.Time `json:"lastLogin"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}
