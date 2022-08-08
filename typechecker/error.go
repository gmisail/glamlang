package typechecker

import "fmt"

type TypeError struct {
	message string
	line    int
}

func CreateTypeError(message string, line int) *TypeError {
	return &TypeError{message, line}
}

func (t *TypeError) Error() string {
	return fmt.Sprintf("[type] line %d, %s\n", t.line, t.message)
}
