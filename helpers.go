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

// isValidPassword checks if a string looks like a valid password
// It should be mostly printable ASCII with minimal replacement characters
func isValidPassword(s string) bool {
	if len(s) < 3 || len(s) > 100 {
		return false
	}

	// Count printable ASCII characters
	printableCount := 0
	replacementCharCount := 0 // Count UTF-8 replacement characters ()

	for _, r := range s {
		if r >= 32 && r < 127 {
			printableCount++
		} else if r == 0xFFFD { // Unicode replacement character
			replacementCharCount++
		}
	}

	// If more than 20% are replacement characters, it's likely binary/garbage
	if replacementCharCount > len(s)/5 {
		return false
	}

	// At least 80% should be printable ASCII
	if printableCount < len(s)*4/5 {
		return false
	}

	// Check for common non-password patterns
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
