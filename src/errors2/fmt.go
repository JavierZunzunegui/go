package errors2

import (
	"fmt"
	"strings"
)

func isError(i interface{}) bool {
	_, ok := i.(error)
	return ok
}

const magicSuffix = ": %w"

// Errorf should remain in package fmt. Placing it here for demo purposes only.
// Mostly copied from go HEAD
func Errorf(format string, a ...interface{}) error {
	frame := Caller(1)

	if len(a) == 0 || !isError(a[len(a)-1]) || !strings.HasSuffix(format, magicSuffix) {
		return &WrappingError{
			frame,
			New(fmt.Sprintf(format, a...)),
			nil,
		}
	}

	err, ok := a[len(a)-1].(error)
	if !ok {
		return &WrappingError{
			frame,
			New(fmt.Sprintf(format, a...)),
			nil,
		}
	}

	return &WrappingError{
		frame,
		New(fmt.Sprintf(format[:len(format)-len(magicSuffix)], a[:len(a)-1]...)),
		toWrappingError(err),
	}
}
