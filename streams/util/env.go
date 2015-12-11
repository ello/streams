package util

import "os"

//GetEnvWithDefault is a convienance method to pull ENV entries with a default value if not present
func GetEnvWithDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

//GetEnvIntWithDefault is a convienance method to pull ENV entries as ints with a default value if not present
func GetEnvIntWithDefault(key string, defaultValue int) int {
	retVal, _ := ValidateInt(os.Getenv(key), defaultValue)

	return retVal
}

//IsEnvPresent will return a boolean of whether the key is present
func IsEnvPresent(key string) bool {
	val := os.Getenv(key)
	return len(val) != 0
}
