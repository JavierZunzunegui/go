package errors2_test

import (
	"errors2"
	"testing"
)

func TestErrorf(t *testing.T) {
	testStringAll(
		t,
		"no_wrapping",
		errors2.Errorf(
			"hello %s",
			"world",
		),
		"hello world",
		"hello world (errors2_test.TestErrorf:wrap_test.go:12)",
		"hello world",
		"hello world"+"\n\t"+"errors2_test.TestErrorf:wrap_test.go:12",
	)

	testStringAll(
		t,
		"single_wrapping",
		errors2.Errorf(
			"I wrap: %w",
			errors2.New("inner error"),
		),
		"I wrap: inner error",
		"I wrap (errors2_test.TestErrorf:wrap_test.go:25): inner error",
		"I wrap"+"\n"+"inner error",
		"I wrap"+"\n\t"+"errors2_test.TestErrorf:wrap_test.go:25"+"\n"+"inner error",
	)
}

func TestWrap(t *testing.T) {
	testStringAll(
		t,
		"nil_wrapping",
		errors2.Wrap(
			nil,
			errors2.New("hello world"),
		),
		"hello world",
		"hello world (errors2_test.TestWrap:wrap_test.go:40)",
		"hello world",
		"hello world"+"\n\t"+"errors2_test.TestWrap:wrap_test.go:40",
	)

	testStringAll(
		t,
		"single_wrapping",
		errors2.Wrap(
			errors2.New("inner error"),
			errors2.New("I wrap"),
		),
		"I wrap: inner error",
		"I wrap (errors2_test.TestWrap:wrap_test.go:53): inner error",
		"I wrap"+"\n"+"inner error",
		"I wrap"+"\n\t"+"errors2_test.TestWrap:wrap_test.go:53"+"\n"+"inner error",
	)

	testStringAll(
		t,
		"custom_error_wraps",
		errors2.Wrap(
			errors2.New("inner error"),
			myError{},
		),
		"custom error type: inner error",
		"custom error type (errors2_test.TestWrap:wrap_test.go:66): inner error",
		"custom error type"+"\n"+"inner error",
		"custom error type"+"\n\t"+"errors2_test.TestWrap:wrap_test.go:66"+"\n"+"inner error",
	)

	testStringAll(
		t,
		"custom_error_wrapped",
		errors2.Wrap(
			myError{},
			errors2.New("I wrap"),
		),
		"I wrap: custom error type",
		"I wrap (errors2_test.TestWrap:wrap_test.go:79): custom error type",
		"I wrap"+"\n"+"custom error type",
		"I wrap"+"\n\t"+"errors2_test.TestWrap:wrap_test.go:79"+"\n"+"custom error type",
	)
}

type myError struct{}

func (myError) Error() string { return "custom error type" }

func testStringAll(
	t *testing.T,
	name string,
	err error,
	expectedDefault, expectedColonLong, expectedMultiLine, expectedMultiLineLong string) {
	colonLongStringer := errors2.NewStringer(
		func() errors2.Formatter { return errors2.NewColonFormatter(true) },
	)
	multiLineStringer := errors2.NewStringer(
		func() errors2.Formatter { return errors2.NewMultiLineFormatter(false) },
	)
	multiLineLongStringer := errors2.NewStringer(
		func() errors2.Formatter { return errors2.NewMultiLineFormatter(true) },
	)

	t.Run(name, func(t *testing.T) {
		testStringSingle(t, "Error_method", err.Error(), expectedDefault)
		testStringSingle(t, "colon_long", colonLongStringer.String(err), expectedColonLong)
		testStringSingle(t, "multi_line", multiLineStringer.String(err), expectedMultiLine)
		testStringSingle(t, "multi_line_long", multiLineLongStringer.String(err), expectedMultiLineLong)
	})
}

func testStringSingle(t *testing.T, name string, got, expected string) {
	t.Run(name, func(t *testing.T) {
		if got != expected {
			t.Fatalf("expected %q got %q", expected, got)
		}
	})
}
