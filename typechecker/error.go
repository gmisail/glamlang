package typechecker

import "fmt"

type TypeError struct {
	message string
}

func CreateTypeError(msg string) *TypeError {
	return &TypeError{message: msg}
}

func (t *TypeError) Error() string {
	return fmt.Sprintf("[type] %s\n", t.message)
}
