package config

import "golang.org/x/oauth2"

func LoadGoogleOauthConfig(id, secret, redirectUrl string) *oauth2.Config {
	var googleOauthConfig = &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		RedirectURL:  redirectUrl,
		Scopes: []string{
			"openid",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}

	return googleOauthConfig
}
