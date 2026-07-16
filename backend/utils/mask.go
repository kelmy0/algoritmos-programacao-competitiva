package utils

import "strings"

func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "[invalid-email-format]"
	}
	local := parts[0]
	domain := parts[1]

	if len(local) <= 2 {
		return "***@" + domain
	}
	return local[:2] + "***@" + domain
}
