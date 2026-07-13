package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

type AuthGoogleHandler struct {
	GoogleConfig        *oauth2.Config
	Service             *services.AuthService
	AppDomain           string
	IsProduce           bool
	RefreshDurationDays int
}

func NewAuthGoogleHandler(service *services.AuthService, googleConfig *oauth2.Config, appDomain string, isProduce bool, refreshDurationDays int) *AuthGoogleHandler {
	return &AuthGoogleHandler{
		Service:             service,
		GoogleConfig:        googleConfig,
		AppDomain:           appDomain,
		IsProduce:           isProduce,
		RefreshDurationDays: refreshDurationDays,
	}
}

func (h *AuthGoogleHandler) GoogleLogin(c *gin.Context) {
	state, err := utils.GenerateCustomId(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error while generating session"})
		return
	}

	c.SetCookie("oauth_google_state", state, 300, "/", h.AppDomain, h.IsProduce, true)

	url := h.GoogleConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthGoogleHandler) GoogleCallback(c *gin.Context) {
	urlState := c.Query("state")
	cookieState, err := c.Cookie("oauth_google_state")
	if err != nil || urlState != cookieState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session expired"})
		return
	}
	c.SetCookie("oauth_google_state", "", -1, "/", h.AppDomain, h.IsProduce, true)

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Google code"})
		return
	}

	token, err := h.GoogleConfig.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	idTokenStr, ok := token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Id Token not found"})
		return
	}

	payload, err := idtoken.Validate(c.Request.Context(), idTokenStr, h.GoogleConfig.ClientID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Google token"})
		return
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is missing from Google token"})
		return
	}

	emailVerified, ok := payload.Claims["email_verified"].(bool)
	if !ok || !emailVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Your Google account email is not verified"})
		return
	}

	socialUserId := payload.Subject

	name, ok := payload.Claims["name"].(string)

	if !ok || name == "" {
		client := h.GoogleConfig.Client(c.Request.Context(), token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err == nil {
			defer resp.Body.Close()
			var googleProfile struct {
				Name string `json:"name"`
			}
			if json.NewDecoder(resp.Body).Decode(&googleProfile) == nil && googleProfile.Name != "" {
				name = googleProfile.Name
			}
		}
	}

	if name == "" {
		name = utils.ExtractNameFromEmail(email)
	}

	result, err := h.Service.AuthWithGoogle(c.Request.Context(), "google", socialUserId, email, name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.RefreshToken != "" {
		c.SetCookie("refresh_token", result.RefreshToken, 60*60*24*h.RefreshDurationDays, "/", h.AppDomain, h.IsProduce, true)
	}
	c.JSON(http.StatusOK, result.LoginResponse)
}
