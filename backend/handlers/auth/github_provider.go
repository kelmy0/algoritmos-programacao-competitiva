package authhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"golang.org/x/oauth2"
)

type GithubProvider struct {
	config *oauth2.Config
}

func NewGithubProvider(config *oauth2.Config) *GithubProvider {
	return &GithubProvider{config: config}
}

func (p *GithubProvider) Name() string                { return "github" }
func (p *GithubProvider) OAuthConfig() *oauth2.Config { return p.config }

func (p *GithubProvider) FetchUser(ctx context.Context, token *oauth2.Token) (*SocialUser, string, error) {
	client := p.config.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		slog.Error("failed to fetch user profile from GitHub API", "error", err)
		return nil, dto.CodeInternalError, nil
	}
	defer resp.Body.Close()

	var ghUser dto.GithubUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		slog.Error("failed to decode GitHub profile response", "error", err)
		return nil, dto.CodeInternalError, nil
	}

	socialUserId := fmt.Sprintf("%d", ghUser.ID)

	emailResp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		slog.Error("failed to fetch user emails from GitHub API", "userId", socialUserId, "error", err)
		return nil, dto.CodeInternalError, nil
	}
	defer emailResp.Body.Close()

	var emails []dto.GithubEmailResponse
	if err := json.NewDecoder(emailResp.Body).Decode(&emails); err != nil {
		slog.Error("failed to decode GitHub emails list", "userId", socialUserId, "error", err)
		return nil, dto.CodeInternalError, nil
	}

	var email string
	for _, e := range emails {
		if e.Primary && e.Verified {
			email = e.Email
			break
		}
	}

	if email == "" {
		return nil, dto.CodeUnverifiedGithubEmail, nil
	}

	name := ghUser.Name
	if name == "" {
		name = ghUser.Login
	}

	return &SocialUser{
		ID:    socialUserId,
		Email: email,
		Name:  name,
	}, "", nil
}
