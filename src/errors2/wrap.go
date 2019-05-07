package errors2

// WrappingError provides the error wrapping functionality, and is exclusively the only error type doing so.
// Each WrappingError holds a (non-WrappingError, non-nil) error as the payload, and points to the next WrappingError.
// The causal error (lowest in the wrapping chain) points to nil.
type WrappingError struct {
	frame   Frame
	payload error
	next    *WrappingError
}

// Error serializes the target according to the default colon serializer, omitting frames and joining with ": ".
// For example, for "wrapper-2" -> "wrapper-1" -> StackError -> "cause" the result would be:
// "wrapper-2: wrapper-1: cause"
func (wErr *WrappingError) Error() string {
	return defaultShortStringer.String(wErr)
}

// Payload is a getter to payload error.
// It is non-nil and never a WrappingError.
func (wErr *WrappingError) Payload() error {
	return wErr.payload
}

// Next is a getter for the next WrappingError in the chain.
// It returns nil for the last error in the chain.
func (wErr *WrappingError) Next() *WrappingError {
	return wErr.next
}

func toWrappingError(err error) *WrappingError {
	wErr, ok := err.(*WrappingError)
	if !ok {
		return &WrappingError{
			Frame{isVoid: true},
			err,
			nil,
		}
	}
	return wErr
}

func SkipWrap(previous, wrap error, skip int) error {
	if wrap == nil {
		// avoid this
		if previous == nil {
			return nil
		}
		return toWrappingError(previous)
	}

	frame := Caller(skip + 1)

	if previous == nil {
		if wErr, ok := wrap.(*WrappingError); ok {
			return wErr
		}

		return &WrappingError{
			frame,
			wrap,
			nil,
		}
	}

	return &WrappingError{
		frame,
		wrap,
		toWrappingError(previous),
	}
}

func Wrap(previous, wrap error) error {
	return SkipWrap(previous, wrap, 1)
}
