package utils

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// RuneReader represents something that can read runes from a bufio.Reader
// but handles \r\n and \n\r as a single '\n'
type RuneReader struct {
	Reader   *bufio.Reader
	position ReaderPosition

	// caching
	hasCache bool // is the cache valid
	cache    rune
}

// NewRuneReaderFromString creates a new RuneReader from a string
func NewRuneReaderFromString(s string) *RuneReader {
	return &RuneReader{
		Reader: bufio.NewReader(strings.NewReader(s)),
	}
}

// NewRuneReaderFromReader creates a new RuneReader from an io.Reader
func NewRuneReaderFromReader(rd io.Reader) *RuneReader {
	return &RuneReader{
		Reader: bufio.NewReader(rd),
	}
}

// Position returns the current position of the reader
// i.e. the position of the next character to be read
func (reader *RuneReader) Position() ReaderPosition {
	return reader.position
}

// Next reads the next character
func (reader *RuneReader) Read() (r rune, pos ReaderPosition, err error) {
	pos.Line = reader.position.Line
	pos.Column = reader.position.Column

	// read the next rune
	a, err := reader.readRaw()
	if err != nil {
		if err == io.EOF {
			pos.EOF = true
			reader.position.EOF = true

			err = nil
		}
		return
	}

	// we consider '\n\r' and '\r\n' special newlines that both get collapsed into just an '\n'
	// since we might have this, we need to have to peek at the next character also
	if a == '\n' || a == '\r' {
		var b rune
		b, err = reader.peekRaw()
		// if we have an EOF next, this will be re-read
		if err == io.EOF {
			r = '\n'
			pos.EOF = false
			err = nil

			// we did have a linebreak
			reader.position.Column = 0
			reader.position.Line++

			return

			// if we have an error, pass that along
		} else if err != nil {
			pos.EOF = false
			return
		}

		// we had a real newline
		if (a == '\n' && b == '\r') || (a == '\r' && b == '\n') {
			reader.eatRaw() // skip the next character
			r = '\n'
			pos.EOF = false
			err = nil

			// we did have a linebreak
			reader.position.Column = 0
			reader.position.Line++

			return
		} else if a == '\n' {
			// we did have a linebreak
			reader.position.Column = 0
			reader.position.Line++
		} else {
			reader.position.Column++
		}

	} else {
		reader.position.Column++
	}

	// return a
	r = a
	return
}

// Eat is like Read, expect that it does not return any values and only an error (if any)
func (reader *RuneReader) Eat() (err error) {
	_, _, err = reader.Read()

	return
}

// Peek peeks into the next character
func (reader *RuneReader) Peek() (r rune, pos ReaderPosition, err error) {
	// read the current position
	pos = reader.Position()

	// peek the next raw character
	r, err = reader.peekRaw()
	if err != nil {
		if err == io.EOF {
			pos.EOF = true
			err = nil
		}
		return
	}

	// if we don't have a potentially special character
	// we can return everything now
	if !(r == '\n' || r == '\r') {
		return
	}

	// else we fallback to read
	r, pos, err = reader.Read()
	if err != nil {
		return
	}

	// and then unread
	err = reader.Unread(r, pos)

	return
}

// Unread unreads a character from the input
func (reader *RuneReader) Unread(r rune, pos ReaderPosition) error {
	if reader.hasCache {
		return errors.New("Cannot unread: Character already cached. ")
	}

	reader.hasCache = true
	reader.cache = r
	reader.position = pos

	return nil
}

// readRaw reads the next character, without taking care of special newlines
// Return err == io.EOF for the end of file
func (reader *RuneReader) readRaw() (r rune, err error) {
	if reader.hasCache {
		// cache and position to return
		r = reader.cache
		reader.hasCache = false

		if reader.position.EOF {
			err = io.EOF
		}

		return
	}

	r, _, err = reader.Reader.ReadRune()
	return
}

// eatRaw eats the next character, without taking care of special newlines
func (reader *RuneReader) eatRaw() (err error) {
	if reader.hasCache {
		reader.hasCache = false
		return
	}
	_, _, err = reader.Reader.ReadRune()
	return
}

// PeekRaw peeks the next character without taking care of special newlines.
func (reader *RuneReader) peekRaw() (r rune, err error) {
	// if we have a cache, return it
	if reader.hasCache {
		r = reader.cache
		return
	}

	// if we do not have this, return
	r, _, err = reader.Reader.ReadRune()
	if err != nil {
		return
	}
	err = reader.Reader.UnreadRune()
	return
}

// ReadWhile reads runes from a string as long as they and the partial string so far match the regexp
// returns the matching string, the first sucessfully read character, and the last successfully read character
// an empty string implies that the first read character did not match, and that loc.Start == loc.End
func (reader *RuneReader) ReadWhile(f func(r rune) bool) (s string, loc ReaderRange, err error) {
	// make a buffer for the string
	var builder strings.Builder
	defer func() {
		if err == nil {
			s = builder.String()
		}
	}()

	// read start position
	loc.Start = reader.position
	loc.End = loc.Start

	// keep reading the current rune
	// as long as there is no EOF
	var r rune
	p := reader.position
	for true {
		loc.End = p
		r, p, err = reader.Read()
		if err != nil {
			return
		}

		// if we reached the end of the string
		// we need to update the position one more time
		if p.EOF {
			loc.End = p
			return
		}

		// if f no longer matches
		// unread and return
		if !f(r) {
			err = reader.Unread(r, p)
			return
		}

		// don't add the trailing EOF
		builder.WriteRune(r)
	}

	return
}

// EatWhile eats runes from the RuneReader as long as f returns true
// Along with a count of how many characters have been eaten
func (reader *RuneReader) EatWhile(f func(r rune) bool) (count int, err error) {
	var r rune
	var p ReaderPosition
	for true {
		r, p, err = reader.Read()
		if err != nil {
			return
		}

		// if we reached EOF, we don't need to check anything
		if p.EOF {
			return
		}

		// if f no longer matches, we unread and return
		if !f(r) {
			err = reader.Unread(r, p)
			return
		}

		count++
	}

	return
}
