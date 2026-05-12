package directory

import (
	"html"
	"regexp"
	"strings"
)

func stripTags(s string) string {
	re := regexp.MustCompile(`(?s)<[^>]*>`)
	return strings.Join(strings.Fields(re.ReplaceAllString(s, " ")), " ")
}

func firstHref(s string) string {
	re := regexp.MustCompile(`(?i)href="([^"]+)"`)
	for _, m := range re.FindAllStringSubmatch(s, -1) {
		if strings.HasPrefix(m[1], "http://") || strings.HasPrefix(m[1], "https://") {
			return html.UnescapeString(m[1])
		}
	}
	return ""
}
