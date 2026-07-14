package dto

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

type ForgotPasswordResponse struct {
	Message string `json:"message"`
}
