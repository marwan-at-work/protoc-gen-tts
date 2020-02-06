package tts

import (
	"fmt"
	"strings"
)

func getDoc(s string, indent int) string {
	if s == "" {
		return ""
	}
	lines := []string{}
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) == 1 {
		return fmt.Sprintf(`/** %s */`, lines[0])
	}
	var resp string
	resp += "/**\n"
	for _, line := range lines {
		resp += strings.Repeat("\t", indent) + line + "\n"
	}
	return resp + strings.Repeat("\t", indent) + "*/"
}
