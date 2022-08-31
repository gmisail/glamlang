package util

import (
	"math"
	"strings"

	"github.com/gmisail/glamlang/io"
	"github.com/gmisail/glamlang/lexer"
)

func Hint(source *io.SourceFile, token *lexer.Token, message string) string {
	var builder strings.Builder

	if token == nil {
		return ""
	}

	if source == nil || token == nil {
		return message
	}

	builder.WriteString(source.GetLine(token.Absolute) + "\n")
	builder.WriteString(strings.Repeat(" ", token.Relative-1))
	builder.WriteString(strings.Repeat("^", token.Length) + "\n")

	hintMidpoint := token.Relative - int(token.Length/2)
	messageMidpoint := len(message) / 2

	// the first half of the message can fit
	builder.WriteString(strings.Repeat(" ", int(math.Max(0, float64(hintMidpoint-messageMidpoint)))))
	builder.WriteString(message)

	return builder.String()
}
