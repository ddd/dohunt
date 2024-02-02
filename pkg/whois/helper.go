package whois

import "unicode"

func isUpper(str string) bool {
	for _, char := range str {
		if !unicode.IsUpper(char) {
			return false
		}
	}
	return true
}
