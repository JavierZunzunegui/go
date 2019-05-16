package errors2_test

import (
	"errors"
	"errors2"
	"fmt"
	"testing"
)

const (
	alternative = "alternative"
	proposal    = "proposal"
)

func benchmarkErrorf(b *testing.B, name string, f func(string, ...interface{}) error) {
	b.Run(name, func(b *testing.B) {
		scenarios := []struct {
			name   string
			format string
			args   []interface{}
		}{
			{
				"no_wrapping_const",
				"hello world",
				[]interface{}{},
			},
			{
				"no_wrapping_formatted",
				"hello %s",
				[]interface{}{"world"},
			},
			{"single_wrapping",
				"I wrap: %w",
				[]interface{}{f("inner error")},
			},
		}

		for _, scenario := range scenarios {
			scenario := scenario

			b.Run(scenario.name, func(b *testing.B) {
				format := scenario.format
				args := scenario.args

				// b.Logf("output(N=%d): (len=%d) %q", b.N, len(f(format, args...).Error()), f(format, args...).Error()) // use for debugging
				b.ResetTimer()

				var err error
				for n := 0; n < b.N; n++ {
					err = f(format, args...)
				}

				// do not optimise the Errorf call away
				_ = err
			})
		}
	})
}

func BenchmarkErrorf(b *testing.B) {
	benchmarkErrorf(b, alternative, errors2.Errorf)
	benchmarkErrorf(b, proposal, fmt.Errorf)
}

func benchmarkErrorMethod(b *testing.B, name string, f func(uint8) error) {
	b.Run(name, func(b *testing.B) {
		scenarios := []struct {
			name  string
			wraps uint8
		}{
			{
				"string_error_len_0",
				0,
			},
			{
				"string_error_len_1",
				1,
			},
			{
				"string_error_len_3",
				3,
			},
			{
				"string_error_len_10",
				10,
			},
			{
				"string_error_len_20",
				20,
			},
		}

		for _, scenario := range scenarios {
			err := f(scenario.wraps)
			b.Run(scenario.name, func(b *testing.B) {
				// b.Logf("output(N=%d): (len=%d) %q", b.N, len(wErr.Error()), wErr.Error()) // use for debugging
				b.ResetTimer()

				var msg string
				for n := 0; n < b.N; n++ {
					msg = err.Error()
				}

				// do not optimise the Error call away
				_ = msg
			})
		}
	})
}

func BenchmarkErrorMethod(b *testing.B) {
	benchmarkErrorMethod(b, alternative, func(wraps uint8) error {
		err := errors2.New("cause")
		for n := 0; n < int(wraps); n++ {
			err = errors2.Errorf("wrapper_%d: %w", n, err)
		}
		return err
	})

	benchmarkErrorMethod(b, proposal, func(wraps uint8) error {
		err := errors.New("cause")
		for n := 0; n < int(wraps); n++ {
			err = fmt.Errorf("wrapper_%d: %w", n, err)
		}
		return err
	})
}
