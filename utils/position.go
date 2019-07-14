package utils

import "fmt"

// ReaderPosition represents the position of a single character within a document
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

// ReaderRange represents a range of read characters from a RuneReader
// The range is defined to contain all characters from Start to End (inclusive).
// A ReaderRange object on its own can not express an empty range, instead the read characters should be returned alongside the range.
type ReaderRange struct {
	Start ReaderPosition
	End   ReaderPosition
}
