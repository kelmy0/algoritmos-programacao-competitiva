package utils

import (
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var HumanNameRegex = regexp.MustCompile(`[^\p{L}\s\-\.\']+`)
var TitleRegex = regexp.MustCompile(`[^\p{L}\s\-\.\'\d\+\#\&]+`)

func SanitizeHumanName(name string) string {
	clean := HumanNameRegex.ReplaceAllString(name, "")
	fields := strings.Fields(clean)
	return strings.Join(fields, " ")
}

func SanitizeTitle(title string) string {
	clean := TitleRegex.ReplaceAllString(title, "")
	fields := strings.Fields(clean)
	return strings.Join(fields, " ")
}

func SanitizeMarkDown(text string) string {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").OnElements("code")
	return p.Sanitize(text)
}

func NormalizeUsername(text string) string {
	re := regexp.MustCompile(`\s+`)
	username := re.ReplaceAllString(text, "_")
	username = strings.Trim(username, "_")
	username = strings.ToLower(username)
	return username
}

func ExtractNameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return "Google User"
}
