package util

import "strconv"

func ToInt(value string) int {
	if value == "" {
		return 0
	}
	if result, err := strconv.Atoi(value); err == nil {
		return result
	}
	return 0
}

func ToBool(value string) bool {
	if value == "" {
		return false
	}
	if result, err := strconv.ParseBool(value); err == nil {
		return result
	}
	return false
}
