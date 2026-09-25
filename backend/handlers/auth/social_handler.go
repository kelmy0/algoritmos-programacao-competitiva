package authhandler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/services/auth"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
	"golang.org/x/oauth2"
)

type socialAuthService interface {
	AuthWithSocialProvider(ctx context.Context, provider, socialUserId, email, name, deviceHash string) (auth.AuthResult, error)
}

type AuthSocialHandler struct {
	providers           map[string]SocialProvider
	service             socialAuthService
	appDomain           string
	isProduce           bool
	refreshDurationDays int
	frontendUrl         string
}

func NewAuthSocialHandler(
	service socialAuthService,
	appDomain, frontendUrl string,
	isProduce bool,
	refreshDurationDays int,
	providers ...SocialProvider,
) *AuthSocialHandler {
	providerMap := make(map[string]SocialProvider)
	for _, p := range providers {
		providerMap[p.Name()] = p
	}

	return &AuthSocialHandler{
		service:             service,
		providers:           providerMap,
		appDomain:           appDomain,
		isProduce:           isProduce,
		refreshDurationDays: refreshDurationDays,
		frontendUrl:         frontendUrl,
	}
}

func (h *AuthSocialHandler) SocialLogin(c *gin.Context) {
	providerName := c.Param("provider")
	provider, ok := h.providers[providerName]
	if !ok {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(dto.CodeInternalError, dto.MsgUnexpectedError))
		return
	}

	stateCookie := fmt.Sprintf("oauth_%s_state", providerName)
	h.startSocialLogin(c, stateCookie, provider.OAuthConfig())
}

func (h *AuthSocialHandler) SocialCallback(c *gin.Context) {
	providerName := c.Param("provider")
	provider, ok := h.providers[providerName]
	if !ok {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(dto.CodeInternalError, dto.MsgUnexpectedError))
		return
	}

	frontendUrl := fmt.Sprintf("%s/auth/callback?", h.frontendUrl)
	stateCookie := fmt.Sprintf("oauth_%s_state", providerName)

	code := h.verifyCallbackAndGetCode(c, stateCookie)
	if code == "" {
		return
	}

	token, err := provider.OAuthConfig().Exchange(c.Request.Context(), code)
	if err != nil {
		c.Redirect(http.StatusFound, frontendUrl+"error="+dto.CodeInternalError)
		return
	}

	socialUser, errorCode, err := provider.FetchUser(c.Request.Context(), token)
	if err != nil || errorCode != "" {
		if errorCode == "" {
			errorCode = dto.CodeInternalError
		}
		c.Redirect(http.StatusFound, frontendUrl+"error="+errorCode)
		return
	}

	deviceHash := ExtractDeviceHash(c.Request)
	result, err := h.service.AuthWithSocialProvider(c.Request.Context(), providerName, socialUser.ID, socialUser.Email, socialUser.Name, deviceHash)
	if err != nil {
		h.socialError(c, err)
		return
	}

	if result.LoginResponse.Requires2FA {
		SetOAuthStateCookie(c, "pre_auth_token", result.LoginResponse.PreAuthToken, h.appDomain, h.isProduce)
		c.Redirect(http.StatusFound, frontendUrl+"pre_auth_token=true")
		return
	}

	if result.RefreshToken != "" {
		SetRefreshCookie(c, result.RefreshToken, h.appDomain, h.refreshDurationDays, h.isProduce)
	}

	if result.LoginResponse.AccessToken != "" {
		SetAccessToken(c, result.LoginResponse.AccessToken, h.appDomain, h.isProduce)
	}

	c.Redirect(http.StatusFound, frontendUrl+"access_token=true")
}

/*
	func (h *AuthSocialHandler) SocialLinkAccount(c *gin.Context) {
		providerName := c.Param("provider")
		provider, ok := h.providers[providerName]
		if !ok {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(dto.CodeInternalError, dto.MsgUnexpectedError))
			return
		}

		linkCookieName := fmt.Sprintf("oauth_%s_link_user", providerName)
		stateCookieName := fmt.Sprintf("oauth_%s_state", providerName)

		h.linkSocialAccount(c, linkCookieName, stateCookieName, provider.OAuthConfig())
	}
*/
func (h *AuthSocialHandler) startSocialLogin(c *gin.Context, cookieName string, config *oauth2.Config) {
	state, err := utils.GenerateCustomId(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(dto.CodeInternalError, dto.MsgUnexpectedError))
		return
	}

	SetOAuthStateCookie(c, cookieName, state, h.appDomain, h.isProduce)

	url := config.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthSocialHandler) verifyCallbackAndGetCode(c *gin.Context, cookieName string) string {
	frontendUrl := fmt.Sprintf("%s/auth/callback?", h.frontendUrl)
	urlState := c.Query("state")
	cookieState, err := c.Cookie(cookieName)
	if err != nil || urlState != cookieState {
		c.Redirect(http.StatusFound, frontendUrl+"error="+dto.CodeSessionExpired)
		return ""
	}
	c.SetCookie(cookieName, "", -1, "/", h.appDomain, h.isProduce, true)

	code := c.Query("code")
	if code == "" {
		c.Redirect(http.StatusFound, frontendUrl+"error="+dto.CodeMissingOAuthCode)
		return ""
	}

	return code
}

/*
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

	c.SetCookie(linkCookieName, userId, 300, "/", h.appDomain, h.isProduce, true)

	h.startSocialLogin(c, cookieName, oauthConfig)
}
*/

func (h *AuthSocialHandler) socialError(c *gin.Context, err error) {
	frontendURL := fmt.Sprintf("%s/auth/callback?error=", h.frontendUrl)

	if appErr, ok := errors.AsType[*models.AppError](err); ok {
		c.Redirect(http.StatusFound, frontendURL+appErr.Code)
		return
	}

	c.Redirect(http.StatusFound, frontendURL+dto.CodeInternalError)
}
