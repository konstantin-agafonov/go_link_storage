// Package e provides error wrapping utilities for consistent error handling.
package e

import "fmt"

// Wrap wraps an error with a message, creating a new error that includes both.
func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

// WrapIfErr wraps an error with a message only if the error is not nil.
// Returns nil if err is nil.
func WrapIfErr(msg string, err error) error {
	if err == nil {
		return nil
	}

	return Wrap(msg, err)
}
