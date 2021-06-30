package utils

import (
	"bufio"
	"io"
	"strings"
)

// RuneReader represents something that can read runes from a RuneReader
// but handles \r\n and \n\r as a single '\n'
type RuneReader struct {
	Reader io.RuneReader

	position ReaderPosition // current position
	pushback []rune         // runes that have been unread
}

// NewRuneReaderFromReader creates a new RuneReader from an io.Reader
func NewRuneReaderFromReader(rd io.Reader) *RuneReader {
	reader := &RuneReader{}

	// if we already have a RuneReader, we can just use it
	// else we make one using bufio!
	if rr, ok := rd.(io.RuneReader); ok {
		reader.Reader = rr
	} else {
		reader.Reader = bufio.NewReader(rd)
	}
	return reader
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

	// we need to take care of special cases '\r\n' and '\n\r'
	// three cases:
	// - read an '\r' => if the next character is an '\n', we need to return it
	// - read an '\n' => the next character doesn't matter, we have an actual '\n'
	// - read anything else => we have the actual characters

	if r != '\r' {
		return
	}

	// read the next character normally, then "unread" it
	r, pos, err = reader.Read()
	if err != nil {
		return
	}
	reader.Unread(r, pos)

	return
}

// Unread unreads a character from the input
func (reader *RuneReader) Unread(r rune, pos ReaderPosition) {
	reader.pushback = append(reader.pushback, r) // store read rune
	reader.position = pos                        // and update position
}

// readRaw reads the next character, without taking care of special newlines
// Return err == io.EOF for the end of file
func (reader *RuneReader) readRaw() (r rune, err error) {
	if len(reader.pushback) > 0 {
		// pop from the pushback
		r, reader.pushback = reader.pushback[len(reader.pushback)-1], reader.pushback[:len(reader.pushback)-1]

		// take care of an EOF
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
	// drop an element from the pushback if we have it
	if len(reader.pushback) > 0 {
		reader.pushback = reader.pushback[:len(reader.pushback)-1]
		return
	}

	_, _, err = reader.Reader.ReadRune()
	return
}

// PeekRaw peeks the next character without taking care of special newlines.
func (reader *RuneReader) peekRaw() (r rune, err error) {
	// if we already have a cached value, just use it!
	if len(reader.pushback) > 0 {
		r = reader.pushback[len(reader.pushback)-1]
		return
	}

	// do a real read, and then unread it!
	r, _, err = reader.Reader.ReadRune()
	if err != nil {
		return
	}
	reader.Unread(r, reader.position) // TODO: We shouldn't have to pass reader.position here!
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

	// keep reading the current rune, until there is an EOF
	var r rune
	p := reader.position
	for {
		loc.End = p
		r, p, err = reader.Read()
		if err != nil {
			return
		}

		// we have reached the end of the stream, so update the position!
		if p.EOF {
			loc.End = p
			return
		}

		// f no longer matches, so we can unread it!
		if !f(r) {
			reader.Unread(r, p)
			return
		}

		// add the read rune to the builder
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
			reader.Unread(r, p)
			return
		}

		count++
	}
}
