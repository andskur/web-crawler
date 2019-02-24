package utils

import (
	"strings"
)

// inSlice checks if slice contain given string
func InSlice(s string, slice []string) bool {
	for _, v := range slice {
		if s == v || v+"/" == s {
			return true
		}
	}
	return false
}

// TODO need refactoring

// inMap check if map contain given link
func InMap(s string, m map[string][]string) bool {
	_, ok := m[s]
	_, okSlash := m[s+"/"]
	_, okOneMore := m[strings.TrimSuffix(s, "/")]
	if okSlash || ok || okOneMore {
		return true
	}
	return false
}
