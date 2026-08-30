package dto

type twoFactorBase struct {
	DeviceHash   string `json:"-"`
	RefreshToken string `json:"-"`
	UserId       string `json:"-"`
}

type TwoFactorGenerateRequest struct {
	Password string `json:"password" binding:"omitempty,min=8"`
	UserId   string `json:"-"`
	Email    string `json:"-"`
}

type TwoFactorGenerateResponse struct {
	Secret string `json:"secret"`
	QRCode string `json:"qrCode"`
}

type TwoFactorEnableRequest struct {
	twoFactorBase
	Code string `json:"code" binding:"required,len=6"`
}

type TwoFactorEnableResponse struct {
	AccessToken string `json:"accessToken"`
}

type TwoFactorDisableRequest struct {
	twoFactorBase
	Password string `json:"password" binding:"omitempty,min=8"`
}
