package utils

import (
	"time"

	"github.com/fatih/structs"
)

// ToDuration converts the string duration to time.Duration, if failed returns the default
func ToDuration(duration string, defaultDuration time.Duration) time.Duration {
	parsedDuration, err := time.ParseDuration(duration)
	if err != nil {
		return defaultDuration
	}
	return parsedDuration
}

// StructToMap converts struct to a map
func StructToMap(data interface{}) map[string]interface{} {
	return structs.Map(data)
}
