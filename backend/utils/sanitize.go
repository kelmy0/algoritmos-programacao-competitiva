package utils

import (
	"html"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var LippinRegex = regexp.MustCompile(`[^\p{L}\s\-\.\d]+`)

func SanitizeName(name string) string {
	cleanName := LippinRegex.ReplaceAllString(name, "")
	fields := strings.Fields(cleanName)
	finalName := strings.Join(fields, " ")
	return html.EscapeString(finalName)
}

func SanitizeMarkDown(text string) string {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").OnElements("code")
	return p.Sanitize(text)
}
