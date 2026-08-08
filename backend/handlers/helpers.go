package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/utils"
)

type CookieConfig struct {
	Name         string
	Value        string
	Domain       string
	Duration     time.Duration
	IsProduction bool
	Partitioned  bool
}

func HandleAPIError(c *gin.Context, err error) {
	if appErr, ok := errors.AsType[*models.AppError](err); ok {
		c.JSON(appErr.StatusCode, dto.NewErrorResponse(appErr.Code, appErr.Message))
		return
	}

	c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
		dto.CodeInternalError,
		dto.MsgUnexpectedError,
	))
}

func GetAuthContext(c *gin.Context) (userId, email, accessId string, accessExpiresAt time.Time, ok bool) {
	userId = c.GetString("userId")
	accessId = c.GetString("accessId")
	accessExpiresAt = c.GetTime("accessExpiresAt")
	email = c.GetString("email")

	if userId == "" || email == "" || accessId == "" || accessExpiresAt.IsZero() {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			dto.CodeMissingUserIdContext,
			dto.MsgMissingDataFromContext,
		))
		c.Abort()
		return "", "", "", time.Time{}, false
	}

	return userId, email, accessId, accessExpiresAt, true
}

func SetRefreshCookie(c *gin.Context, value, domain string, refreshDurationDays int, isProduction bool) {
	cookie := buildSecureCookie(CookieConfig{
		Name:         "refresh_token",
		Value:        value,
		Domain:       domain,
		Duration:     time.Duration(refreshDurationDays) * 24 * time.Hour,
		IsProduction: isProduction,
		Partitioned:  true,
	})

	c.SetCookieData(cookie)
}

func SetAccessToken(c *gin.Context, value, domain string, isProduction bool) {
	cookie := buildSecureCookie(CookieConfig{
		Name:         "access_token",
		Value:        value,
		Domain:       domain,
		Duration:     15 * time.Minute,
		IsProduction: isProduction,
		Partitioned:  true,
	})
	c.SetCookieData(cookie)
}

func SetOAuthStateCookie(c *gin.Context, name, value, domain string, isProduction bool) {
	cookie := buildSecureCookie(CookieConfig{
		Name:         name,
		Value:        value,
		Domain:       domain,
		Duration:     5 * time.Minute,
		IsProduction: isProduction,
		Partitioned:  false,
	})

	c.SetCookieData(cookie)
}

func ClearCookie(c *gin.Context, name, domain string, isProduction bool) {
	cookie := buildSecureCookie(CookieConfig{
		Name:         name,
		Value:        "",
		Domain:       domain,
		Duration:     -1 * time.Hour,
		IsProduction: isProduction,
	})

	c.SetCookieData(cookie)
}

func ExtractDeviceHash(r *http.Request) string {
	fmt.Printf("DEBUG -> UA: %q | Lang: %q | Plat: %q | Mob: %q\n",
		r.UserAgent(),
		r.Header.Get("Accept-Language"),
		r.Header.Get("Sec-CH-UA-Platform"),
		r.Header.Get("Sec-CH-UA-Mobile"),
	)

	return utils.GenerateDeviceHash(
		r.UserAgent(),
		r.Header.Get("Accept-Language"),
		r.Header.Get("Sec-CH-UA-Platform"),
		r.Header.Get("Sec-CH-UA-Mobile"))
}

func buildSecureCookie(cfg CookieConfig) *http.Cookie {
	seconds := int(cfg.Duration.Seconds())

	cookie := &http.Cookie{
		Name:     cfg.Name,
		Value:    cfg.Value,
		Path:     "/",
		Domain:   cfg.Domain,
		MaxAge:   seconds,
		Secure:   cfg.IsProduction,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	if cfg.Duration > 0 {
		cookie.Expires = time.Now().Add(cfg.Duration)
	}

	if cfg.Partitioned && cfg.IsProduction {
		cookie.Unparsed = append(cookie.Unparsed, "Partitioned")
	}

	return cookie
}
