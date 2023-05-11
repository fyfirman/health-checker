package helpers

import (
	"fmt"
	"strconv"
	"strings"
)

func ReplaceWildcard(input string, numDigits int, replacementNumber int) string {
	// Find the position of the wildcard character
	wildcardPos := strings.Index(input, "*")
	if wildcardPos == -1 {
		// Wildcard not found, return original string
		return input
	}

	// Generate the replacement string with the specified number of digits
	formatStr := "%0" + strconv.Itoa(numDigits) + "d"
	replacement := fmt.Sprintf(formatStr, replacementNumber)

	// Replace the wildcard character with the generated replacement string
	output := strings.Replace(input, "*", replacement, 1)
	return output
}
