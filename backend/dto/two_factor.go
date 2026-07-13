package dto

type TwoFactorGenerateResponse struct {
	Secret string `json:"secret"`
	QRCode string `json:"qr_code"`
}

type TwoFactorEnableRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

type TwoFactorDisableRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

type TwoFactorEnableResponse struct {
	Message string `json:"message"`
}
