package authhandler

import (
	"context"

	"golang.org/x/oauth2"
)

type SocialUser struct {
	ID    string
	Email string
	Name  string
}

type SocialProvider interface {
	Name() string
	OAuthConfig() *oauth2.Config
	FetchUser(ctx context.Context, token *oauth2.Token) (*SocialUser, string, error)
}
