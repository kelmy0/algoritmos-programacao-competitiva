package utils

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var HumanNameRegex = regexp.MustCompile(`[^\p{L}\s\-\.\']+`)
var TitleRegex = regexp.MustCompile(`[^\p{L}\s\-\.\'\d\+\#\&]+`)
var UsernameRegex = regexp.MustCompile(`[^\p{L}\p{N}_\-]+`)
var MultipleSpacesRegex = regexp.MustCompile(`\s+`)

func SanitizeHumanName(name string) string {
	clean := HumanNameRegex.ReplaceAllString(name, "")
	fields := strings.Fields(clean)
	return cases.Title(language.BrazilianPortuguese).String(strings.ToLower(strings.Join(fields, " ")))
}

func SanitizeUsername(username string) string {
	clean := UsernameRegex.ReplaceAllString(username, "")
	fields := strings.Fields(clean)
	return strings.ToLower(strings.Join(fields, ""))
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
	clean := MultipleSpacesRegex.ReplaceAllString(text, "_")
	clean = UsernameRegex.ReplaceAllString(clean, "")
	if clean == "" {
		clean = "user"
	}
	id, _ := GenerateCustomId(12)

	maxCleanLen := 32 - 1 - len(id)
	if len(clean) > maxCleanLen {
		clean = clean[:maxCleanLen]
	}

	return clean + "_" + id
}

func ExtractNameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 && parts[0] != "" {
		cleanPart := strings.ReplaceAll(parts[0], "+", " ")
		cleanPart = strings.ReplaceAll(cleanPart, ".", " ")

		formattedName := SanitizeHumanName(cleanPart)

		if utf8.RuneCountInString(formattedName) >= 6 {
			return formattedName
		}
	}
	return "Social User"
}
