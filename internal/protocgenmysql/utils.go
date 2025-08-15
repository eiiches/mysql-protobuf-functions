package protocgenmysql

func EscapeSQLString(s string) string {
	// Escape single quotes and backslashes for SQL string literals
	result := ""
	for _, char := range s {
		switch char {
		case '\'':
			result += "''"
		case '\\':
			result += "\\\\"
		default:
			result += string(char)
		}
	}
	return result
}
