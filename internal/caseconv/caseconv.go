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
