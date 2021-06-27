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

// NewRuneReaderFromReader creates a new RuneReader from an io.Reader
func NewRuneReaderFromReader(rd io.Reader) *RuneReader {
	return &RuneReader{
		Reader: bufio.NewReader(rd),
	}
}

// NewRuneReaderFromString creates a new RuneReader from a string
func NewRuneReaderFromString(s string) *RuneReader {
	return NewRuneReaderFromReader(strings.NewReader(s))
}

// Position returns the current position of the reader
// i.e. the position of the next character to be read
func (reader *RuneReader) Position() ReaderPosition {
	return reader.position
}

// Read reads the next character from the input and returns it
func (reader *RuneReader) Read() (r rune, pos ReaderPosition, err error) {
	// tell the caller that we read the current position
	pos.Line = reader.position.Line
	pos.Column = reader.position.Column

	// read the next rune, and handle errors!
	r, err = reader.readRaw()
	if err != nil {
		if err == io.EOF {
			pos.EOF = true
			reader.position.EOF = true

			err = nil
		}
		return
	}

	// handle '\r\n' and '\n\r' as a special newline and normalize them into a single '\n'
	if r == '\n' || r == '\r' {
		// lookahead to the next character
		var l rune
		l, err = reader.peekRaw()

		// handle errors in the lookahead
		if err != nil {
			if err == io.EOF {
				// io.EOF will be re-read in
				err = nil

				// we did have a linebreak
				reader.position.Column = 0
				reader.position.Line++
			}
			pos.EOF = false
			return
		}

		// we caught a collapsed newline, and should collapse both into a single line.
		if (r == '\n' && l == '\r') || (r == '\r' && l == '\n') {
			reader.eatRaw() // skip the next character
			r = '\n'
			pos.EOF = false
			err = nil
		}
	}

	// increment the position of the reader!
	if r != '\n' {
		reader.position.Column++
	} else {
		reader.position.Column = 0
		reader.position.Line++
	}

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
	pos = reader.position

	// peek the next raw character
	r, err = reader.peekRaw()
	if err != nil {
		if err == io.EOF {
			pos.EOF = true
			err = nil
		}
		return
	}

	// we need to take care of special characters: '\r\n' and '\n\r'
	// if we have a '\n', the next character is guaranteed to be a '\n'
	// (even though the character after might be skipped)
	if r != '\r' {
		return
	}

	// however, if we have an '\r', we need to look at the next character too
	// but as we are peeking, we need to cache that
	reader.cache, reader.position, err = reader.Read()
	if err != nil {
		return
	}
	reader.hasCache = true

	// return the cached value
	r = reader.cache
	pos = reader.position
	return
}

// ErrNoUnread is returned when a character cannot be unread
var ErrNoUnread = errors.New("Unread: Already has a cached character")

// Unread unreads a character from the input
func (reader *RuneReader) Unread(r rune, pos ReaderPosition) error {
	if reader.hasCache {
		return ErrNoUnread
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

	// keep reading the current rune
	// as long as there is no EOF
	var r rune
	p := reader.position
	for {
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
}

// EatWhile eats runes from the RuneReader as long as f returns true
// Along with a count of how many characters have been eaten
func (reader *RuneReader) EatWhile(f func(r rune) bool) (count int, err error) {
	var r rune
	var p ReaderPosition
	for {
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
}
