package executil

import (
	"bytes"
	"os/exec"
	"strings"
)

func FormatCommand(cmd *exec.Cmd) (string, error) {
	var out bytes.Buffer
	for i, arg := range cmd.Args {
		if i > 0 {
			out.WriteRune(' ')
		}
		out.WriteString(QuoteWord(arg))
	}
	return out.String(), nil
}

func QuoteWord(word string) string {
	if strings.ContainsAny(word, `'" `) {
		return `"` + strings.Replace(word, `"`, `\"`, -1) + `"`
	} else {
		return word
	}
}
