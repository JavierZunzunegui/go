package errors2

import (
	"bytes"
)

// Formatter defines how errors will be converted to a string.
// They should be rare - most users will never have to create one.
// They are stateful.
type Formatter interface {
	Init(*WrappingError)

	Next() (error, Frame)

	Format(error, Frame, *bytes.Buffer)
}

type inOrderPartialFormatter struct {
	// final
	formatFrames bool

	// variable
	firstEntry bool
	currentErr *WrappingError
}

func (f *inOrderPartialFormatter) Init(wErr *WrappingError) {
	f.currentErr = wErr
	f.firstEntry = true
}

func (f *inOrderPartialFormatter) Next() (error, Frame) {
	if f.currentErr == nil {
		return nil, Frame{}
	}
	payload := f.currentErr.payload
	frame := f.currentErr.frame
	f.currentErr = f.currentErr.next
	return payload, frame
}

type colonFormatter struct {
	inOrderPartialFormatter
}

func (s *colonFormatter) Format(err error, frame Frame, buf *bytes.Buffer) {
	if s.firstEntry {
		s.firstEntry = false
	} else {
		buf.WriteString(": ")
	}

	if bufErr, ok := err.(BufferError); ok {
		bufErr.ErrorToBuffer(buf)
	} else {
		buf.WriteString(err.Error())
	}

	if s.formatFrames && !frame.isVoid {
		buf.WriteString(" (")
		frame.Format(buf)
		buf.WriteString(")")
	}
}

// NewColonFormatter provides a formatter that appends messages with ': '.
// It is the Formatter used by the %s representation of errors.
func NewColonFormatter(formatFrames bool) Formatter {
	return &colonFormatter{inOrderPartialFormatter: inOrderPartialFormatter{formatFrames: formatFrames}}
}

var (
	defaultShortStringer = NewStringer(func() Formatter { return NewColonFormatter(false) })
)

type multiLineFormatter struct {
	inOrderPartialFormatter
}

func (s *multiLineFormatter) Format(err error, frame Frame, buf *bytes.Buffer) {
	if s.firstEntry {
		s.firstEntry = false
	} else {
		buf.WriteString("\n")
	}

	if bufErr, ok := err.(BufferError); ok {
		bufErr.ErrorToBuffer(buf)
	} else {
		buf.WriteString(err.Error())
	}

	if s.formatFrames && !frame.isVoid {
		buf.WriteString("\n\t")
		frame.Format(buf)
	}
}

// NewColonFormatter provides a formatter that appends messages with '\n'.
func NewMultiLineFormatter(formatFrames bool) Formatter {
	return &multiLineFormatter{inOrderPartialFormatter: inOrderPartialFormatter{formatFrames: formatFrames}}
}
