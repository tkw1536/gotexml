package utils

import "fmt"

// ReaderPosition represents the position of the reader
type ReaderPosition struct {
	// current 0-based line
	Line uint

	// current 0-based column
	Column uint

	// true iff we are at the end of the input
	EOF bool
}

// String turns this ReaderPosition into a string
func (rp ReaderPosition) String() string {
	s := ""
	if rp.EOF {
		s = " (at EOF)"
	}
	return fmt.Sprintf("line %d column %d%s", rp.Line, rp.Column, s)
}

// ReaderRange represents a range within a read document
type ReaderRange struct {
	Start ReaderPosition
	End   ReaderPosition
}
