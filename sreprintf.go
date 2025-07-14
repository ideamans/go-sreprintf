package sreprintf

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Sreprintf extracts values from a formatted message using a template,
// then applies those values to a translation template.
func Sreprintf(template, message, translation string) (string, error) {
	// Extract values from the message using the template
	values, err := extractValues(template, message)
	if err != nil {
		return "", err
	}

	// Apply extracted values to the translation template
	return applyValues(translation, values)
}

// placeholderInfo stores information about a placeholder
type placeholderInfo struct {
	start        int
	end          int
	placeholder  string
	verb         byte
	hasFlag      bool
	hasWidth     bool
	hasPrecision bool
}

// parsePlaceholder parses a placeholder and extracts its components
func parsePlaceholder(s string) placeholderInfo {
	info := placeholderInfo{placeholder: s}

	if len(s) < 2 || s[0] != '%' {
		return info
	}

	// Extract verb (last character)
	info.verb = s[len(s)-1]

	// Check for flags, width, precision
	middle := s[1 : len(s)-1]
	if strings.ContainsAny(middle, "#+-") {
		info.hasFlag = true
	}
	if strings.ContainsAny(middle, "0123456789") {
		info.hasWidth = true
	}
	if strings.Contains(middle, ".") {
		info.hasPrecision = true
	}

	return info
}

// extractValues extracts values from a formatted message using a template
func extractValues(template, message string) ([]interface{}, error) {
	var values []interface{}

	// Find all placeholders in template
	re := regexp.MustCompile(`%(?:[#+-])?(?:[0-9]+)?(?:\.[0-9]+)?[sdftv%xXobBeEfFgGqp]`)
	matches := re.FindAllStringIndex(template, -1)

	if len(matches) == 0 {
		// No placeholders, check if template matches message exactly
		if template == message {
			return values, nil
		}
		return nil, fmt.Errorf("message does not match template")
	}

	// Parse placeholders
	var placeholders []placeholderInfo
	for _, match := range matches {
		ph := template[match[0]:match[1]]
		if ph != "%%" {
			info := parsePlaceholder(ph)
			info.start = match[0]
			info.end = match[1]
			placeholders = append(placeholders, info)
		}
	}

	// Build pattern to extract values
	pattern := buildExtractionPattern(template, matches)

	// Compile and match
	patternRe, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile pattern: %w", err)
	}

	submatches := patternRe.FindStringSubmatch(message)
	if submatches == nil {
		return nil, fmt.Errorf("message does not match template")
	}

	// Convert captured strings to appropriate types
	for i, info := range placeholders {
		if i+1 >= len(submatches) {
			break
		}

		value := submatches[i+1]
		converted := convertValue(value, info.verb)
		values = append(values, converted)
	}

	return values, nil
}

// buildExtractionPattern builds a regex pattern to extract values from formatted string
func buildExtractionPattern(template string, placeholderMatches [][]int) string {
	var pattern strings.Builder
	pattern.WriteString("^")

	pos := 0
	for _, match := range placeholderMatches {
		// Add literal text before placeholder
		if pos < match[0] {
			pattern.WriteString(regexp.QuoteMeta(template[pos:match[0]]))
		}

		placeholder := template[match[0]:match[1]]

		// Handle literal percent
		if placeholder == "%%" {
			pattern.WriteString("%")
		} else {
			// Add capture group based on placeholder type
			info := parsePlaceholder(placeholder)
			pattern.WriteString(getCapturePattern(info))
		}

		pos = match[1]
	}

	// Add remaining literal text
	if pos < len(template) {
		pattern.WriteString(regexp.QuoteMeta(template[pos:]))
	}

	pattern.WriteString("$")
	return pattern.String()
}

// getCapturePattern returns appropriate regex capture pattern for a placeholder
func getCapturePattern(info placeholderInfo) string {
	switch info.verb {
	case 'd':
		if info.hasWidth || info.hasFlag {
			// For width-specified or flag-specified integers, be more flexible
			return `(\s*[+-]?\d+)`
		}
		// Allow non-numeric values for type mismatch tolerance
		return `([^\s]+)`

	case 'f', 'F', 'e', 'E', 'g', 'G':
		if info.hasWidth || info.hasFlag {
			// For formatted floats, be more flexible with whitespace and signs
			return `(\s*[+-]?\d+(?:\.\d+)?)`
		}
		return `(-?\d+(?:\.\d+)?)`

	case 's', 'v':
		if info.hasWidth {
			// For width-specified strings, match including spaces
			return `(.+?)`
		}
		// For regular strings, match non-greedy any character
		return `(.*?)`

	case 't':
		return `(true|false)`

	case 'x', 'X':
		return `((?:0[xX])?[0-9a-fA-F]+)`

	case 'b', 'o':
		return `(\d+)`

	default:
		return `(.+?)`
	}
}

// convertValue converts a string value to the appropriate type based on verb
func convertValue(value string, verb byte) interface{} {
	value = strings.TrimSpace(value)

	switch verb {
	case 'd', 'b', 'o':
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
		return value

	case 'f', 'F', 'e', 'E', 'g', 'G':
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			return v
		}
		return value

	case 't':
		if v, err := strconv.ParseBool(value); err == nil {
			return v
		}
		return value

	case 'x', 'X':
		// Try to parse as hex integer
		cleanValue := strings.TrimPrefix(strings.TrimPrefix(value, "0x"), "0X")
		if v, err := strconv.ParseInt(cleanValue, 16, 64); err == nil {
			return int(v)
		}
		return value

	default:
		return value
	}
}

// applyValues applies extracted values to a translation template
func applyValues(translation string, values []interface{}) (string, error) {
	// Count placeholders in translation (excluding %%)
	re := regexp.MustCompile(`%(?:[#+-])?(?:[0-9]+)?(?:\.[0-9]+)?[sdftv%xXobBeEfFgGqp]`)
	matches := re.FindAllString(translation, -1)

	placeholderCount := 0
	for _, match := range matches {
		if match != "%%" {
			placeholderCount++
		}
	}

	// Adjust values array to match placeholder count
	if len(values) > placeholderCount {
		values = values[:placeholderCount]
	}
	for len(values) < placeholderCount {
		values = append(values, "")
	}

	// Apply values using fmt.Sprintf
	result := fmt.Sprintf(translation, values...)
	return result, nil
}

