package helper

import "strings"

func MakeSlug(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	replacer := strings.NewReplacer(
		" ", "-",
		"/", "-",
		",", "",
		".", "",
		":", "",
	)

	return replacer.Replace(value)
}
