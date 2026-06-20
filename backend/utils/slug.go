package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var slugCleanupRegex = regexp.MustCompile(`[^a-z0-6\-]+`)

func Slug(name string) string {
	str := strings.ToLower(name)

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	strWithNoAccents, _, _ := transform.String(t, str)

	strWithHiphons := strings.ReplaceAll(strWithNoAccents, " ", "-")
	finalSlug := slugCleanupRegex.ReplaceAllString(strWithHiphons, "")

	var multipleHyphensRegex = regexp.MustCompile(`-+`)
	finalSlug = multipleHyphensRegex.ReplaceAllString(finalSlug, "-")
	return strings.Trim(finalSlug, "-")
}
