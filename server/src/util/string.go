package util

import "strings"

func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

func IsNotBlank(s string) bool {
	return !IsBlank(s)
}
