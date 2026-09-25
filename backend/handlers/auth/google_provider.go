package authhandler

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
	"golang.org/x/oauth2"
)

type GoogleProvider struct {
	config *oauth2.Config
}

func NewGoogleProvider(config *oauth2.Config) *GoogleProvider {
	return &GoogleProvider{config: config}
}

func (p *GoogleProvider) Name() string                { return "google" }
func (p *GoogleProvider) OAuthConfig() *oauth2.Config { return p.config }

func (p *GoogleProvider) FetchUser(ctx context.Context, token *oauth2.Token) (*SocialUser, string, error) {
	idTokenStr, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, dto.CodeMissingTokenID, nil
	}

	payload, err := idtoken.Validate(ctx, idTokenStr, p.config.ClientID)
	if err != nil {
		return nil, dto.CodeInvalidGoogleToken, nil
	}

	googleUser, errorCode := dto.NewGoogleUserPayload(payload)
	if errorCode != "" {
		return nil, errorCode, nil
	}

	name := googleUser.Name
	if name == "" {
		name = p.fetchNameFallback(ctx, token)
		if name == "" {
			name = utils.ExtractNameFromEmail(googleUser.Email)
		}
	}

	return &SocialUser{
		ID:    googleUser.Subject,
		Email: googleUser.Email,
		Name:  name,
	}, "", nil
}

func (p *GoogleProvider) fetchNameFallback(ctx context.Context, token *oauth2.Token) string {
	client := p.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var googleProfile struct {
		Name string `json:"name"`
	}
	if json.NewDecoder(resp.Body).Decode(&googleProfile) == nil {
		return googleProfile.Name
	}
	return ""

}
