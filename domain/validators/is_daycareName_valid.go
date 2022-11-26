package validators

import "strings"

func IsDaycareNameValid(name string) bool {
	newName := name

	if strings.Contains(newName, name) {
		return false
	}

	return true
}
