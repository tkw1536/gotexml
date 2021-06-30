package utils

import (
	"fmt"

	"github.com/pkg/errors"
)

// ReaderError represents an error of an error
type ReaderError struct {
	error
	Reader   *RuneReader
	Location ReaderPosition
}

// Cause returns the underlying cause of this error
func (r ReaderError) Cause() error {
	return r.error
}

// Error returns the error
func (r *ReaderError) Error() string {
	return fmt.Sprintf("%s near %s", r.error.Error(), r.Location)
}

// NewErrorF returns a new error for the given reader and error message
func NewErrorF(reader *RuneReader, format string, args ...interface{}) *ReaderError {
	return &ReaderError{
		error:    fmt.Errorf(format, args...),
		Reader:   reader,
		Location: reader.Position(),
	}
}

// WrapErrorF returns a new error for the given reader and underlying error
func WrapErrorF(reader *RuneReader, err error, format string, args ...interface{}) *ReaderError {
	return &ReaderError{
		error:    errors.Wrapf(err, format, args...),
		Reader:   reader,
		Location: reader.Position(),
	}
}
