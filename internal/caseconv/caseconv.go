package caseconv

import "strings"

func SnakeToUpperCamel(s string) string {
	builder := strings.Builder{}
	nextUpper := true
	for _, ch := range strings.ToLower(s) {
		if ch == '_' {
			nextUpper = true
			continue
		}
		if nextUpper && 'a' <= ch && ch <= 'z' {
			builder.WriteRune(ch - 'a' + 'A')
		} else {
			builder.WriteRune(ch)
		}
		nextUpper = false
	}
	return builder.String()
}

func LowerCamelToSnake(s string) string {
	if s == "" {
		return ""
	}

	builder := strings.Builder{}
	for i, ch := range s {
		if 'A' <= ch && ch <= 'Z' {
			// Add underscore before uppercase letters (except the first character)
			if i > 0 {
				builder.WriteRune('_')
			}
			// Convert to lowercase
			builder.WriteRune(ch - 'A' + 'a')
		} else {
			builder.WriteRune(ch)
		}
	}
	return builder.String()
}
