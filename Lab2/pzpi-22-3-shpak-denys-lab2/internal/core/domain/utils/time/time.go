// Package time provides utility functions for time operations.
package time // import "wayra/internal/core/domain/utils/time"

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseDuration parses a duration string in the format "HH:MM:SS" and returns a time.Duration.
// durationStr: a string in the format "HH:MM:SS".
// Returns a time.Duration and an error if the duration string is invalid.
func ParseDuration(durationStr string) (time.Duration, error) {
	parts := strings.Split(durationStr, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid duration format: %s", durationStr)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hour part: %v", err)
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minute part: %v", err)
	}
	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf("invalid second part: %v", err)
	}

	totalDuration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
	return totalDuration, nil
}
