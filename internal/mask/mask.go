package mask

import "regexp"

// panRegex matches 16-digit PANs with optional spaces/dashes.
var panRegex = regexp.MustCompile(`\b(?:\d{4}[- ]?){3}\d{4}\b`)

// maskPAN replaces the numeric PAN with “**** **** **** ” + last 4 digits.
func maskPAN(s string) string {
	if len(s) < 4 {
		return s
	}
	return "**** **** **** " + s[len(s)-4:]
}
