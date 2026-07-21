package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
	"golang.org/x/oauth2"
)

type AuthSocialHandler struct {
	GoogleConfig        *oauth2.Config
	GithubConfig        *oauth2.Config
	Service             *services.AuthService
	AppDomain           string
	IsProduce           bool
	RefreshDurationDays int
	FrontendUrl         string
}

func NewAuthSocialHandler(service *services.AuthService, googleConfig, githubConfig *oauth2.Config, appDomain, frontendUrl string, isProduce bool, refreshDurationDays int) *AuthSocialHandler {
	return &AuthSocialHandler{
		Service:             service,
		GithubConfig:        githubConfig,
		GoogleConfig:        googleConfig,
		AppDomain:           appDomain,
		IsProduce:           isProduce,
		RefreshDurationDays: refreshDurationDays,
		FrontendUrl:         frontendUrl,
	}
}

func (h *AuthSocialHandler) GoogleLogin(c *gin.Context) {
	h.startSocialLogin(c, "oauth_google_state", h.GoogleConfig)
}

func (h *AuthSocialHandler) GoogleCallback(c *gin.Context) {
	code := h.verifyCallbackAndGetCode(c, "oauth_google_state")
	if code == "" {
		return
	}

	frontendUrl := fmt.Sprintf("%s/auth/callback?", h.FrontendUrl)
	token, err := h.GoogleConfig.Exchange(c.Request.Context(), code)
	if err != nil {
		c.Redirect(http.StatusFound, frontendUrl+"error="+dto.CodeInternalError)
		return
	}

	idTokenStr, ok := token.Extra("id_token").(string)
	if !ok {
		c.Redirect(http.StatusFound, frontendUrl+"error="+dto.CodeMissingTokenID)
		return
	}

	payload, err := idtoken.Validate(c.Request.Context(), idTokenStr, h.GoogleConfig.ClientID)
	if err != nil {
		c.Redirect(http.StatusFound, frontendUrl+"error="+dto.CodeInvalidGoogleToken)
		return
	}

	googleUser, errorCode := dto.NewGoogleUserPayload(payload)
	if errorCode != "" {
		c.Redirect(http.StatusFound, frontendUrl+"error="+errorCode)
		return
	}

	if googleUser.Name == "" {
		googleUser.Name = h.fetchGoogleNameFallback(c.Request.Context(), token)
		if googleUser.Name == "" {
			googleUser.Name = utils.ExtractNameFromEmail(googleUser.Email)
		}
	}
	linkedUserID, err := c.Cookie("oauth_google_link_user")

	if err == nil && linkedUserID != "" {
		c.SetCookie("oauth_google_link_user", "", -1, "/", h.AppDomain, h.IsProduce, true)

		err := h.Service.LinkSocialAccount(c.Request.Context(), linkedUserID, "google", googleUser.Subject, googleUser.Email)
		if err != nil {
			h.socialError(c, err)
			return
		}

		c.Redirect(http.StatusFound, frontendUrl)
		return
	}

	result, err := h.Service.AuthWithSocialProvider(c.Request.Context(), "google", googleUser.Subject, googleUser.Email, googleUser.Name)
	if err != nil {
		h.socialError(c, err)
		return
	}

	if result.LoginResponse.Requires2FA {
		frontendURL := fmt.Sprintf("%s/auth/callback?pre_token=%s", h.FrontendUrl, result.LoginResponse.PreAuthToken)
		c.Redirect(http.StatusFound, frontendURL)
		return
	}

	if result.RefreshToken != "" {
		c.SetCookie("refresh_token", result.RefreshToken, 60*60*24*h.RefreshDurationDays, "/", h.AppDomain, h.IsProduce, true)
	}

	frontendURL := fmt.Sprintf("%s/auth/callback?access_token=%s", h.FrontendUrl, result.LoginResponse.AccessToken)
	c.Redirect(http.StatusFound, frontendURL)
}

func (h *AuthSocialHandler) GoogleLinkAccount(c *gin.Context) {
	h.linkSocialAccount(c, "oauth_google_link_user", "oauth_google_state", h.GoogleConfig)
}

func (h *AuthSocialHandler) GithubLogin(c *gin.Context) {
	h.startSocialLogin(c, "oauth_github_state", h.GithubConfig)
}

func (h *AuthSocialHandler) GithubCallback(c *gin.Context) {
	code := h.verifyCallbackAndGetCode(c, "oauth_github_state")
	if code == "" {
		return
	}

	frontendUrl := fmt.Sprintf("%s/auth/callback?", h.FrontendUrl)
	token, err := h.GithubConfig.Exchange(c.Request.Context(), code)
	if err != nil {
		slog.Error("failed to exchange code for GitHub token", "error", err)
		c.Redirect(http.StatusFound, frontendUrl+"error="+dto.CodeInternalError)
		return
	}

	client := h.GithubConfig.Client(c.Request.Context(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		slog.Error("failed to fetch user profile from GitHub API", "error", err)
		c.Redirect(http.StatusFound, frontendUrl+"error="+dto.CodeInternalError)
		return
	}
	defer resp.Body.Close()

	var ghUser dto.GithubUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		slog.Error("failed to decode GitHub profile response", "error", err)
		c.Redirect(http.StatusFound, frontendUrl+"error="+dto.CodeInternalError)
		return
	}

	socialUserId := fmt.Sprintf("%d", ghUser.ID)

	emailResp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		slog.Error("failed to fetch user emails from GitHub API", "userId", socialUserId, "error", err)
		c.Redirect(http.StatusFound, frontendUrl+"error="+dto.CodeInternalError)
		return
	}
	defer emailResp.Body.Close()

	var emails []dto.GithubEmailResponse
	if err := json.NewDecoder(emailResp.Body).Decode(&emails); err != nil {
		slog.Error("failed to decode GitHub emails list", "userId", socialUserId, "error", err)
		c.Redirect(http.StatusFound, frontendUrl+"error="+dto.CodeInternalError)
		return
	}

	var email string
	for _, e := range emails {
		if e.Primary && e.Verified {
			email = e.Email
			break
		}
	}

	if email == "" {
		c.Redirect(http.StatusFound, frontendUrl+"error="+dto.CodeUnverifiedGithubEmail)
		return
	}

	name := ghUser.Name
	if name == "" {
		name = ghUser.Login
	}

	linkedUserID, err := c.Cookie("oauth_github_link_user")
	if err == nil && linkedUserID != "" {
		c.SetCookie("oauth_github_link_user", "", -1, "/", h.AppDomain, h.IsProduce, true)

		err := h.Service.LinkSocialAccount(c.Request.Context(), linkedUserID, "github", socialUserId, email)
		if err != nil {
			h.socialError(c, err)
			return
		}

		c.Redirect(http.StatusFound, frontendUrl)
		return
	}

	result, err := h.Service.AuthWithSocialProvider(c.Request.Context(), "github", socialUserId, email, name)
	if err != nil {
		h.socialError(c, err)
		return
	}

	if result.RefreshToken != "" {
		c.SetCookie("refresh_token", result.RefreshToken, 60*60*24*h.RefreshDurationDays, "/", h.AppDomain, h.IsProduce, true)
	}

	frontendURL := fmt.Sprintf("%s/auth/callback?access_token=%s", h.FrontendUrl, result.LoginResponse.AccessToken)
	c.Redirect(http.StatusFound, frontendURL)
}

func (h *AuthSocialHandler) GithubLinkAccount(c *gin.Context) {
	h.linkSocialAccount(c, "oauth_github_link_user", "oauth_github_state", h.GithubConfig)
}

func (h *AuthSocialHandler) startSocialLogin(c *gin.Context, cookieName string, config *oauth2.Config) {
	state, err := utils.GenerateCustomId(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(dto.CodeInternalError, dto.MsgUnexpectedError))
		return
	}

	c.SetCookie(cookieName, state, 300, "/", h.AppDomain, h.IsProduce, true)

	url := config.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthSocialHandler) verifyCallbackAndGetCode(c *gin.Context, cookieName string) string {
	frontendUrl := fmt.Sprintf("%s/auth/callback?", h.FrontendUrl)
	urlState := c.Query("state")
	cookieState, err := c.Cookie(cookieName)
	if err != nil || urlState != cookieState {
		c.Redirect(http.StatusFound, frontendUrl+"error="+dto.CodeSessionExpired)
		return ""
	}
	c.SetCookie(cookieName, "", -1, "/", h.AppDomain, h.IsProduce, true)

	code := c.Query("code")
	if code == "" {
		c.Redirect(http.StatusFound, frontendUrl+"error="+dto.CodeMissingOAuthCode)
		return ""
	}

	return code
}

func (h *AuthSocialHandler) fetchGoogleNameFallback(ctx context.Context, token *oauth2.Token) string {
	client := h.GoogleConfig.Client(ctx, token)
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

func (h *AuthSocialHandler) linkSocialAccount(c *gin.Context, linkCookieName, cookieName string, oauthConfig *oauth2.Config) {
	rawUserId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(dto.CodeInternalError, dto.MsgUnexpectedError))
		return
	}

	userId, ok := rawUserId.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(dto.CodeInternalError, dto.MsgUnexpectedError))
		return
	}

	c.SetCookie(linkCookieName, userId, 300, "/", h.AppDomain, h.IsProduce, true)

	h.startSocialLogin(c, cookieName, oauthConfig)
}

func (h *AuthSocialHandler) socialError(c *gin.Context, err error) {
	frontendURL := fmt.Sprintf("%s/auth/callback?error=", h.FrontendUrl)

	if appErr, ok := errors.AsType[*models.AppError](err); ok {

		c.Redirect(http.StatusFound, frontendURL+appErr.Code)
		return
	}

	c.Redirect(http.StatusFound, frontendURL+dto.CodeInternalError)
}
