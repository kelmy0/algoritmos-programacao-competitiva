package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

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
		Endpoint: google.Endpoint,
	}

	return googleOauthConfig
}

func LoadGithubOauthConfig(id, secret, redirectUrl string) *oauth2.Config {
	var githubOauthConfig = &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		RedirectURL:  redirectUrl,
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     github.Endpoint,
	}

	return githubOauthConfig
}
