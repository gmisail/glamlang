package typechecker

import "fmt"

type TypeError struct {
	message string
	line    int
}

func CreateTypeError(msg string) *TypeError {
	return &TypeError{message: msg, line: 0}
}

func (t *TypeError) Error() string {
	return fmt.Sprintf("[type] %s\n", t.message)
}
