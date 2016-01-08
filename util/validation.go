package util

import "strconv"

// ValidateInt validates that the value is an int (or empty string) or returns the default value
func ValidateInt(value string, defaultVal int) (int, error) {
	if value == "" {
		return defaultVal, nil
	}
	parsedVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultVal, err
	}
	return parsedVal, nil
}
