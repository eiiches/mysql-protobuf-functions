package dedent

import "strings"

func Pipe(s string) string {
	builder := strings.Builder{}

	endsWithNewline := false

	lines := strings.Split(s, "\n")

	for i, line := range lines {
		line = strings.TrimLeft(line, " \t")
		if i == 0 && line == "" {
			continue
		}
		if i == len(lines)-1 && line == "" {
			endsWithNewline = true
			break
		}
		if !strings.HasPrefix(line, "|") {
			panic("dedent.Pipe: line does not start with '|'")
		}
		builder.WriteString(line[1:]) // trim leading pipe
		builder.WriteString("\n")
	}

	if endsWithNewline {
		return builder.String()
	} else {
		result := builder.String()
		if result == "" {
			return ""
		}
		return result[0 : len(result)-1]
	}
}
