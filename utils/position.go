package utils

// ReaderPosition represents the position of the reader
type ReaderPosition struct {
	// current 0-based line
	Line uint

	// current 0-based column
	Column uint

	// true iff we are at the end of the input
	EOF bool
}
