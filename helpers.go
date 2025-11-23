package main

import (
	"strings"
	"unicode"
)

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func removeNonPrintableAscii(input string) string {
	var resultBuilder []rune

	for _, char := range input {
		if unicode.IsPrint(char) && char >= 32 && char != 127 {
			resultBuilder = append(resultBuilder, char)
		}
	}

	return string(resultBuilder)
}

func isValidPassword(s string) bool {
	if len(s) < 3 || len(s) > 100 {
		return false
	}

	printableCount := 0
	replacementCharCount := 0

	for _, r := range s {
		if r >= 32 && r < 127 {
			printableCount++
		} else if r == 0xFFFD {
			replacementCharCount++
		}
	}

	if replacementCharCount > len(s)/5 {
		return false
	}

	if printableCount < len(s)*4/5 {
		return false
	}

	lower := strings.ToLower(s)
	if strings.HasPrefix(lower, "fsha256") ||
		strings.HasPrefix(lower, "ssh-") ||
		strings.HasPrefix(lower, "ecdsa-") ||
		strings.HasPrefix(lower, "rsa-") ||
		strings.HasPrefix(lower, "ed25519") ||
		strings.Contains(lower, "curve25519") ||
		strings.Contains(lower, "diffie-hellman") ||
		strings.Contains(lower, "sntrup") {
		return false
	}

	return true
}
