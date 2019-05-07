package errors2

import (
	"bytes"
	"runtime"
	"strconv"
	"strings"
)

// A Frame contains part of a call stack.
type Frame struct {
	isVoid bool
	frames [1]uintptr
}

// Caller returns a Frame that describes a frame on the caller's stack.
// The argument skip is the number of frames to skip over.
// Caller(0) returns the frame for the caller of Caller.
func Caller(skip int) Frame {
	var s Frame
	runtime.Callers(skip+2, s.frames[:])
	return s
}

// location reports the file, line, and function of a frame.
//
// The returned function may be "" even if file and line are not.
func (f Frame) data() (function, file string, line int) {
	frames := runtime.CallersFrames(f.frames[:])
	fr, _ := frames.Next()
	return fr.Function, fr.File, fr.Line
}

func shortened(s string) string {
	if i := strings.LastIndex(s, "/"); i != -1 {
		return s[i+1:]
	}
	return s
}

func (f Frame) Format(buf *bytes.Buffer) {
	function, file, line := f.data()
	if function != "" {
		buf.WriteString(shortened(function))
	}
	if function != "" && file != "" {
		buf.WriteString(":")
	}
	if file != "" {
		buf.WriteString(shortened(file))
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(line))
	}
}
