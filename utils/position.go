package utils

import "fmt"

// ReaderPosition represents the position of a character within a document
type ReaderPosition struct {
	Line   uint `json:"line"`          // zero-based line number
	Column uint `json:"column"`        // zero-based column number
	EOF    bool `json:"eof,omitempty"` // true iff we are at the end of the input
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
	Start ReaderPosition `json:"start"` // the first character included within the range
	End   ReaderPosition `json:"end"`   // the last character included within the range
}
