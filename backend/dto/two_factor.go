package dto

type requirePassword struct {
	Password string `json:"password" binding:"required,min=8"`
}

type TwoFactorGenerateRequest struct {
	requirePassword
}

type TwoFactorGenerateResponse struct {
	Secret string `json:"secret"`
	QRCode string `json:"qrCode"`
}

type TwoFactorEnableRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

type TwoFactorDisableRequest struct {
	requirePassword
}
