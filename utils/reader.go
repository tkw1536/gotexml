package utils

import (
	"bufio"
	"errors"
	"io"
	"regexp"
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

// GetPosition returns the current position of the reader
func (raw *RuneReader) GetPosition() ReaderPosition {
	return raw.position
}

// Next reads the next character
func (raw *RuneReader) Read() (r rune, pos ReaderPosition, err error) {
	// catch the end-of line
	var eof bool
	defer func() {
		// if something went wrong, don't change the state
		if err != nil {
			return
		}

		// return the current position
		pos = ReaderPosition{
			Line:   raw.position.Line,
			Column: raw.position.Column,
			EOF:    eof,
		}

		// and update the position

		if eof {
			// we are at the end => store eof and don't change position
			raw.position.EOF = true

		} else if r == '\n' {
			// we have a newline => set column to 0 and increase line counter
			raw.position.Column = 0
			raw.position.Line++

		} else {
			// else we only increase the column
			raw.position.Column++
		}
	}()

	// read the next rune
	a, err := raw.readRaw()
	if err != nil {
		if err == io.EOF {
			eof = true
			err = nil
		}
		return
	}

	// we consider '\n\r' and '\r\n' special newlines that both get collapsed into just an '\n'
	// since we might have this, we need to have to peek at the next character also
	if a == '\n' || a == '\r' {
		var b rune
		b, err = raw.peekRaw()
		// if we have an EOF next, this will be re-read
		if err == io.EOF {
			r = '\n'
			eof = false
			err = nil
			return

			// if we have an error, pass that along
		} else if err != nil {
			eof = false
			return
		}

		// we had a real newline
		if (a == '\n' && b == '\r') || (a == '\r' && b == '\n') {
			raw.eatRaw() // skip the next character
			r = '\n'
			eof = false
			err = nil
			return
		}
	}

	// return a
	r = a
	return
}

// Eat is like Read, expect that it does not return any values and only an error (if any)
func (raw *RuneReader) Eat() (err error) {
	_, _, err = raw.Read()

	return
}

// Peek peeks into the next character
func (raw *RuneReader) Peek() (r rune, pos ReaderPosition, err error) {
	// read the current position
	pos = raw.GetPosition()

	// peek the next raw character
	r, err = raw.peekRaw()
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
	r, pos, err = raw.Read()
	if err != nil {
		return
	}

	// and then unread
	err = raw.unreadRaw(r, pos, false)

	return
}

// Unread unreads a character
func (raw *RuneReader) Unread(r rune, pos ReaderPosition) error {
	return raw.unreadRaw(r, pos, false)
}

// readRaw reads the next character, without taking care of special newlines
// Return err == io.EOF for the end of file
func (raw *RuneReader) readRaw() (r rune, err error) {
	if raw.hasCache {
		// cache and position to return
		r = raw.cache
		raw.hasCache = false

		if raw.position.EOF {
			err = io.EOF
		}

		return
	}

	r, _, err = raw.Reader.ReadRune()
	return
}

// eatRaw eats the next character, without taking care of special newlines
func (raw *RuneReader) eatRaw() (err error) {
	if raw.hasCache {
		raw.hasCache = false
		return
	}
	_, _, err = raw.Reader.ReadRune()
	return
}

// PeekRaw peeks the next character without taking care of special newlines.
func (raw *RuneReader) peekRaw() (r rune, err error) {
	// if we have a cache, return it
	if raw.hasCache {
		r = raw.cache
		return
	}

	// if we do not have this, return
	r, _, err = raw.Reader.ReadRune()
	if err != nil {
		return
	}
	err = raw.Reader.UnreadRune()
	return
}

// unreadRaw unreads a character
// When unsafe is true, the check for already having a cached character is skipped
func (raw *RuneReader) unreadRaw(r rune, pos ReaderPosition, unsafe bool) (err error) {
	if raw.hasCache && !unsafe {
		return errors.New("Cannot unread: Character already cached. ")
	}

	raw.hasCache = true
	raw.cache = r
	raw.position = pos

	return
}

// ReadWhile reads runes from a string as long as they and the partial string so far match the regexp
// returns the matching string, the position of the first character of the string, and the last position of the string
// an empty string implies that the first read character did not match, and that loc.Start == loc.End
func (raw *RuneReader) ReadWhile(f func(r rune) bool) (s string, loc ReaderRange, err error) {
	// make a buffer for the string
	var builder strings.Builder
	defer func() {
		if err == nil {
			s = builder.String()
		}
	}()

	// read start position
	loc.Start = raw.position
	loc.End = loc.Start

	// keep reading the current rune
	// as long as there is no EOF
	var r rune
	var p ReaderPosition
	for true {
		r, p, err = raw.Read()
		if err != nil {
			return
		}

		// if we reached EOF, we don't need to check anything
		if p.EOF {
			loc.End = p
			return
		}

		// if f no longer matches
		// unread and return
		if !f(r) {
			err = raw.Unread(r, p)
			return
		}

		// don't add the trailing EOF
		builder.WriteRune(r)
		loc.End = p
	}

	return
}

// ReadWhileMatch reads runes as long as re matches the characters read up to this point
// When not in an error state, it is guaranteed that either re.MatchString(s) is true or s is empty and start === end
func (raw *RuneReader) ReadWhileMatch(re *regexp.Regexp) (s string, loc ReaderRange, err error) {
	var builder strings.Builder
	return raw.ReadWhile(func(r rune) bool {
		builder.WriteRune(r)
		return re.MatchString(builder.String())
	})
}

// EatWhile eats runes from the RuneReader as long as f returns true
// The first rune for which f does not match is not eaten
func (raw *RuneReader) EatWhile(f func(r rune) bool) (count int, err error) {
	var r rune
	var p ReaderPosition
	for true {
		r, p, err = raw.Read()
		if err != nil {
			return
		}

		// if we reached EOF, we don't need to check anything
		if p.EOF {
			return
		}

		// if f no longer matches, we unread and return
		if !f(r) {
			err = raw.Unread(r, p)
			return
		}

		count++
	}

	return
}

// EatWhileMatch eats characters from RuneReader as long as re matches the substring so far
// the first substring for which the substring does not match is not eaten
func (raw *RuneReader) EatWhileMatch(re *regexp.Regexp) (count int, err error) {
	s, _, err := raw.ReadWhileMatch(re)
	count = len(s)
	return
}
