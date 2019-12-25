package tfdocgen

import "strings"

func indefiniteArticle(s string) string {
	vowels := []string{"a", "e", "i", "o", "u", "A", "E", "I", "O", "U"}
	for _, c := range vowels {
		if strings.HasSuffix(s, c) {
			return "an"
		}
	}
	return "a"
}
