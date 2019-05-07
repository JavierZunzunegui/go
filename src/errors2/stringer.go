package errors2

import (
	"bytes"
	"sync"
)

func NewStringer(fFactory func() Formatter) *Stringer {
	return &Stringer{
		pool: sync.Pool{
			New: func() interface{} {
				return &stringerPoolData{
					f: fFactory(),
				}
			},
		},
	}
}

type stringerPoolData struct {
	f   Formatter
	buf bytes.Buffer
}

// Stringer is a thin wrapper around a Formatter factory that converts any wrapped error to its string format.
type Stringer struct {
	pool sync.Pool
}

func (s *Stringer) String(err error) string {
	pd := s.pool.Get().(*stringerPoolData)

	s.toBuffer(pd.f, err, &pd.buf)

	out := pd.buf.String()

	pd.buf.Reset()
	s.pool.Put(pd)

	return out
}

func (s *Stringer) toBuffer(f Formatter, err error, buf *bytes.Buffer) {
	wErr, ok := err.(*WrappingError)
	if !ok {
		wErr = &WrappingError{
			payload: err,
		}
	}

	f.Init(wErr)

	for err, frame := f.Next(); err != nil; err, frame = f.Next() {
		f.Format(err, frame, buf)
	}
}

// BufferError is an optional interface that errors may implement for efficiency purposes.
// If an error's Error() method results in a string allocation for the return statement, it would benefit from this.
type BufferError interface {
	// ErrorToBuffer is a buffering equivalent to Error().
	// The same string as is returned by Error() is to be written to the buffer.
	ErrorToBuffer(*bytes.Buffer)
}
