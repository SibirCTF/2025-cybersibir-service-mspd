package config

import (
	"strings"
)

var KeyDict [10]string

func padWithZeros(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat("0", length-len(s))
}

func init() {
	// TODO: use random key generation
	KeyDict = [10]string{"nerisande", "neriael", "neriyuko", "nerielys", "nerysgosa", "nerieth", "neriett", "neridana", "neriss", "neri"}
	for i := range KeyDict {
		KeyDict[i] = padWithZeros(KeyDict[i], 32)
	}
}
