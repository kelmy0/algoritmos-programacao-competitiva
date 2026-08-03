package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type TurnstileResponse struct {
	Success     bool     `json:"success"`
	ErrorCodes  []string `json:"error-codes"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
}

func VerifyTurnstile(secretKey, token, remoteIP string) (bool, error) {
	if token == "" {
		return false, fmt.Errorf("token is empty")
	}

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify", url.Values{
		"secret":   {secretKey},
		"response": {token},
		"remoteip": {remoteIP},
	})

	if err != nil {
		return false, fmt.Errorf("failed to connect to Cloudflare: %w", err)
	}
	defer resp.Body.Close()

	var result TurnstileResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to decode Cloudflare responsee: %w", err)
	}

	return result.Success, nil
}
