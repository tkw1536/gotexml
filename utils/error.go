package utils

import (
	"fmt"

	"github.com/pkg/errors"
)

// ReaderError represents an error of an error
type ReaderError struct {
	Reader   *RuneReader
	Location ReaderPosition

	inner error
}

// Error returns the error
func (r *ReaderError) Error() string {
	return fmt.Sprintf("%s near %s", r.inner.Error(), r.Location)
}

// NewErrorF returns a new error for the given reader and error message
func NewErrorF(reader *RuneReader, format string, args ...interface{}) *ReaderError {
	return &ReaderError{
		Reader:   reader,
		Location: reader.Position(),
		inner:    fmt.Errorf(format, args...),
	}
}

// WrapErrorF returns a new error for the given reader and underlying error
func WrapErrorF(reader *RuneReader, err error, format string, args ...interface{}) *ReaderError {
	return &ReaderError{
		Reader:   reader,
		Location: reader.Position(),
		inner:    errors.Wrapf(err, format, args...),
	}
}
