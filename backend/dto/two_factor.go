package dto

type TwoFactorGenerateRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

type TwoFactorGenerateResponse struct {
	Secret string `json:"secret"`
	QRCode string `json:"qrCode"`
}

type TwoFactorEnableRequest struct {
	Code         string `json:"code" binding:"required,len=6"`
	DeviceHash   string `json:"-"`
	RefreshToken string `json:"-"`
	UserId       string `json:"-"`
}

type TwoFactorEnableResponse struct {
	AccessToken string `json:"accessToken"`
}

type TwoFactorDisableRequest struct {
	TwoFactorGenerateRequest
	DeviceHash   string `json:"-"`
	RefreshToken string `json:"-"`
	UserId       string `json:"-"`
}
