package shared

import "unicode"

func CheckAlphanumeric(str string) bool {
	var hasLetter, hasNumber bool

	for _, letter := range str {
		switch {
		case unicode.IsLetter(letter):
			hasLetter = true
		case unicode.IsNumber(letter):
			hasNumber = true
		}

		if hasLetter && hasNumber {
			return true
		}
	}

	return false
}
