package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

const testInput = "line 1\nline 2\r\nline 3\n\rline 4\nline 5\ralso\r\rstill\n\nline 7\n"

// makeTestReader makes a new test reader
func makeTestReader() *RuneReader {
	return NewRuneReaderFromString(testInput)
}

// the testOutput we are expecting
var testOutput = []struct {
	wantR   rune
	wantPos ReaderPosition
}{
	{'l', ReaderPosition{0, 0, false}},
	{'i', ReaderPosition{0, 1, false}},
	{'n', ReaderPosition{0, 2, false}},
	{'e', ReaderPosition{0, 3, false}},
	{' ', ReaderPosition{0, 4, false}},
	{'1', ReaderPosition{0, 5, false}},
	{'\n', ReaderPosition{0, 6, false}},

	{'l', ReaderPosition{1, 0, false}},
	{'i', ReaderPosition{1, 1, false}},
	{'n', ReaderPosition{1, 2, false}},
	{'e', ReaderPosition{1, 3, false}},
	{' ', ReaderPosition{1, 4, false}},
	{'2', ReaderPosition{1, 5, false}},
	{'\n', ReaderPosition{1, 6, false}},

	{'l', ReaderPosition{2, 0, false}},
	{'i', ReaderPosition{2, 1, false}},
	{'n', ReaderPosition{2, 2, false}},
	{'e', ReaderPosition{2, 3, false}},
	{' ', ReaderPosition{2, 4, false}},
	{'3', ReaderPosition{2, 5, false}},
	{'\n', ReaderPosition{2, 6, false}},

	{'l', ReaderPosition{3, 0, false}},
	{'i', ReaderPosition{3, 1, false}},
	{'n', ReaderPosition{3, 2, false}},
	{'e', ReaderPosition{3, 3, false}},
	{' ', ReaderPosition{3, 4, false}},
	{'4', ReaderPosition{3, 5, false}},
	{'\n', ReaderPosition{3, 6, false}},

	{'l', ReaderPosition{4, 0, false}},
	{'i', ReaderPosition{4, 1, false}},
	{'n', ReaderPosition{4, 2, false}},
	{'e', ReaderPosition{4, 3, false}},
	{' ', ReaderPosition{4, 4, false}},
	{'5', ReaderPosition{4, 5, false}},
	{'\r', ReaderPosition{4, 6, false}},

	{'a', ReaderPosition{4, 7, false}},
	{'l', ReaderPosition{4, 8, false}},
	{'s', ReaderPosition{4, 9, false}},
	{'o', ReaderPosition{4, 10, false}},

	{'\r', ReaderPosition{4, 11, false}},
	{'\r', ReaderPosition{4, 12, false}},

	{'s', ReaderPosition{4, 13, false}},
	{'t', ReaderPosition{4, 14, false}},
	{'i', ReaderPosition{4, 15, false}},
	{'l', ReaderPosition{4, 16, false}},
	{'l', ReaderPosition{4, 17, false}},

	{'\n', ReaderPosition{4, 18, false}},
	{'\n', ReaderPosition{5, 0, false}},

	{'l', ReaderPosition{6, 0, false}},
	{'i', ReaderPosition{6, 1, false}},
	{'n', ReaderPosition{6, 2, false}},
	{'e', ReaderPosition{6, 3, false}},
	{' ', ReaderPosition{6, 4, false}},
	{'7', ReaderPosition{6, 5, false}},

	{'\n', ReaderPosition{6, 6, false}},

	{rune(0), ReaderPosition{7, 0, true}},
	{rune(0), ReaderPosition{7, 0, true}},
}

func TestRuneReader_Read(t *testing.T) {

	raw := makeTestReader()

	for i, tt := range testOutput {
		t.Run(fmt.Sprintf("read character %d", i), func(t *testing.T) {
			// read the current position
			gotPos := raw.GetPosition()
			gotPos.EOF = tt.wantPos.EOF // ignore EOF (as this may be different)
			if !reflect.DeepEqual(gotPos, tt.wantPos) {
				t.Errorf("RuneReader.GetPosition() gotPOS = %v, want %v", gotPos, tt.wantPos)
			}

			gotR, gotPos, err := raw.Read()
			if (err != nil) != false {
				t.Errorf("RuneReader.Read() error = %v, wantErr %v", err, false)
				return
			}
			if gotR != tt.wantR {
				t.Errorf("RuneReader.Read() gotR = %v, want %v", gotR, tt.wantR)
			}
			if !reflect.DeepEqual(gotPos, tt.wantPos) {
				t.Errorf("RuneReader.Read() gotPOS = %v, want %v", gotPos, tt.wantPos)
			}
		})
	}
}

func TestRuneReader_PeekEat(t *testing.T) {

	raw := makeTestReader()

	for i, tt := range testOutput {
		t.Run(fmt.Sprintf("peek + eat character %d", i), func(t *testing.T) {

			// read the current position
			gotPos := raw.GetPosition()
			gotPos.EOF = tt.wantPos.EOF // ignore EOF (as this may be different)
			if !reflect.DeepEqual(gotPos, tt.wantPos) {
				t.Errorf("RuneReader.GetPosition() gotPOS = %v, want %v", gotPos, tt.wantPos)
			}

			// peek the next character
			gotR, gotPos, err := raw.Peek()
			if (err != nil) != false {
				t.Errorf("RuneReader.Peek() error = %v, wantErr %v", err, false)
				return
			}
			if gotR != tt.wantR {
				t.Errorf("RuneReader.Peek() gotR = %v, want %v", gotR, tt.wantR)
			}
			if !reflect.DeepEqual(gotPos, tt.wantPos) {
				t.Errorf("RuneReader.Peek() gotPOS = %v, want %v", gotPos, tt.wantPos)
			}

			// read the current position (again)
			gotPos = raw.GetPosition()
			gotPos.EOF = tt.wantPos.EOF // ignore EOF (as this may be different)
			if !reflect.DeepEqual(gotPos, tt.wantPos) {
				t.Errorf("RuneReader.GetPosition() gotPOS = %v, want %v", gotPos, tt.wantPos)
			}

			// eat
			err = raw.Eat()
			if (err != nil) != false {
				t.Errorf("RuneReader.Eat() error = %v, wantErr %v", err, false)
				return
			}

		})
	}
}

func TestRuneReader_ReadUnread(t *testing.T) {

	raw := makeTestReader()

	for i, tt := range testOutput {
		t.Run(fmt.Sprintf("read + unread + reread character %d", i), func(t *testing.T) {

			// read the current position
			gotPos := raw.GetPosition()
			gotPos.EOF = tt.wantPos.EOF // ignore EOF (as this may be different)
			if !reflect.DeepEqual(gotPos, tt.wantPos) {
				t.Errorf("RuneReader.GetPosition() gotPOS = %v, want %v", gotPos, tt.wantPos)
			}

			// read the next character
			gotR, gotPos, err := raw.Read()
			if (err != nil) != false {
				t.Errorf("RuneReader.Read() error = %v, wantErr %v", err, false)
				return
			}
			if gotR != tt.wantR {
				t.Errorf("RuneReader.Read() gotR = %v, want %v", gotR, tt.wantR)
			}
			if !reflect.DeepEqual(gotPos, tt.wantPos) {
				t.Errorf("RuneReader.Read() gotPOS = %v, want %v", gotPos, tt.wantPos)
			}

			// unread it
			err = raw.Unread(gotR, gotPos)
			if (err != nil) != false {
				t.Errorf("RuneReader.Unread() error = %v, wantErr %v", err, false)
				return
			}

			// read the current position (again)
			gotPos = raw.GetPosition()
			gotPos.EOF = tt.wantPos.EOF // ignore EOF (as this may be different)
			if !reflect.DeepEqual(gotPos, tt.wantPos) {
				t.Errorf("RuneReader.GetPosition() gotPOS = %v, want %v", gotPos, tt.wantPos)
			}

			// re-read the next character
			gotR, gotPos, err = raw.Read()
			if (err != nil) != false {
				t.Errorf("RuneReader.Read() error = %v, wantErr %v", err, false)
				return
			}
			if gotR != tt.wantR {
				t.Errorf("RuneReader.Read() gotR = %v, want %v", gotR, tt.wantR)
			}
			if !reflect.DeepEqual(gotPos, tt.wantPos) {
				t.Errorf("RuneReader.Read() gotPOS = %v, want %v", gotPos, tt.wantPos)
			}

		})
	}
}

func TestRuneReader_ReadUnreadPeekEat(t *testing.T) {

	raw := makeTestReader()

	for i, tt := range testOutput {
		t.Run(fmt.Sprintf("read + unread + peek + eat character %d", i), func(t *testing.T) {
			// read the current position
			gotPos := raw.GetPosition()
			gotPos.EOF = tt.wantPos.EOF // ignore EOF (as this may be different)
			if !reflect.DeepEqual(gotPos, tt.wantPos) {
				t.Errorf("RuneReader.GetPosition() gotPOS = %v, want %v", gotPos, tt.wantPos)
			}

			// read the next character
			gotR, gotPos, err := raw.Read()
			if (err != nil) != false {
				t.Errorf("RuneReader.Read() error = %v, wantErr %v", err, false)
				return
			}
			if gotR != tt.wantR {
				t.Errorf("RuneReader.Read() gotR = %v, want %v", gotR, tt.wantR)
			}
			if !reflect.DeepEqual(gotPos, tt.wantPos) {
				t.Errorf("RuneReader.Read() gotPOS = %v, want %v", gotPos, tt.wantPos)
			}

			// unread it
			err = raw.Unread(gotR, gotPos)
			if (err != nil) != false {
				t.Errorf("RuneReader.Unread() error = %v, wantErr %v", err, false)
				return
			}

			// read the current position (again)
			gotPos = raw.GetPosition()
			gotPos.EOF = tt.wantPos.EOF // ignore EOF (as this may be different)
			if !reflect.DeepEqual(gotPos, tt.wantPos) {
				t.Errorf("RuneReader.GetPosition() gotPOS = %v, want %v", gotPos, tt.wantPos)
			}

			// re-peak the next character
			gotR, gotPos, err = raw.Peek()
			if (err != nil) != false {
				t.Errorf("RuneReader.Peek() error = %v, wantErr %v", err, false)
				return
			}
			if gotR != tt.wantR {
				t.Errorf("RuneReader.Peek() gotR = %v, want %v", gotR, tt.wantR)
			}
			if !reflect.DeepEqual(gotPos, tt.wantPos) {
				t.Errorf("RuneReader.Peek() gotPOS = %v, want %v", gotPos, tt.wantPos)
			}

			// read the current position (again again)
			gotPos = raw.GetPosition()
			gotPos.EOF = tt.wantPos.EOF // ignore EOF (as this may be different)
			if !reflect.DeepEqual(gotPos, tt.wantPos) {
				t.Errorf("RuneReader.GetPosition() gotPOS = %v, want %v", gotPos, tt.wantPos)
			}

			// and eat it
			err = raw.Eat()
			if (err != nil) != false {
				t.Errorf("RuneReader.Eat() error = %v, wantErr %v", err, false)
				return
			}
		})
	}
}

func TestRuneReader_ReadWhile(t *testing.T) {
	tests := []struct {
		name    string
		f       func(r rune) bool
		wantS   string
		wantLoc ReaderRange
	}{
		{
			"read nothing",
			func(r rune) bool { return false },
			"",
			ReaderRange{
				ReaderPosition{0, 0, false},
				ReaderPosition{0, 0, false},
			},
		},
		{
			"read letters",
			func(r rune) bool { return (r >= 'a' && r <= 'z') },
			"line",
			ReaderRange{
				ReaderPosition{0, 0, false},
				ReaderPosition{0, 3, false},
			},
		},
		{
			"read until the end",
			func(r rune) bool { return true },
			"line 1\nline 2\nline 3\nline 4\nline 5\ralso\r\rstill\n\nline 7\n",
			ReaderRange{
				ReaderPosition{0, 0, false},
				ReaderPosition{7, 0, true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := makeTestReader()
			gotS, gotLoc, err := raw.ReadWhile(tt.f)
			if (err != nil) != false {
				t.Errorf("RuneReader.ReadWhile() error = %v, wantErr %v", err, false)
				return
			}
			if gotS != tt.wantS {
				t.Errorf("RuneReader.ReadWhile() gotS = %v, want %v", gotS, tt.wantS)
			}
			if !reflect.DeepEqual(gotLoc, tt.wantLoc) {
				t.Errorf("RuneReader.ReadWhile() gotLoc = %v, want %v", gotLoc, tt.wantLoc)
			}
		})
	}
}

func TestRuneReader_ReadWhileMatch(t *testing.T) {
	tests := []struct {
		name    string
		re      *regexp.Regexp
		wantS   string
		wantLoc ReaderRange
	}{
		{
			"read nothing",
			regexp.MustCompile("^$"),
			"",
			ReaderRange{
				ReaderPosition{0, 0, false},
				ReaderPosition{0, 0, false},
			},
		},
		{
			"read letters",
			regexp.MustCompile("^[aA-zZ]+$"),
			"line",
			ReaderRange{
				ReaderPosition{0, 0, false},
				ReaderPosition{0, 3, false},
			},
		},
		{
			"read until the end",
			regexp.MustCompile(".*"),
			"line 1\nline 2\nline 3\nline 4\nline 5\ralso\r\rstill\n\nline 7\n",
			ReaderRange{
				ReaderPosition{0, 0, false},
				ReaderPosition{7, 0, true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := makeTestReader()

			gotS, gotLoc, err := raw.ReadWhileMatch(tt.re)
			if (err != nil) != false {
				t.Errorf("RuneReader.ReadWhileMatch() error = %v, wantErr %v", err, false)
				return
			}
			if gotS != tt.wantS {
				t.Errorf("RuneReader.ReadWhileMatch() gotS = %v, want %v", gotS, tt.wantS)
			}
			if !reflect.DeepEqual(gotLoc, tt.wantLoc) {
				t.Errorf("RuneReader.ReadWhileMatch() gotLoc = %v, want %v", gotLoc, tt.wantLoc)
			}
		})
	}
}

func TestRuneReader_EatWhile(t *testing.T) {
	tests := []struct {
		name      string
		f         func(r rune) bool
		wantCount int
		wantPos   ReaderPosition
	}{
		{
			"eat nothing",
			func(r rune) bool { return false },
			0,
			ReaderPosition{0, 0, false},
		},
		{
			"eat letters",
			func(r rune) bool { return (r >= 'a' && r <= 'z') },
			4,
			ReaderPosition{0, 4, false},
		},
		{
			"eat everything",
			func(r rune) bool { return true },
			55,
			ReaderPosition{7, 0, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := makeTestReader()
			gotCount, gotErr := raw.EatWhile(tt.f)
			if (gotErr != nil) != false {
				t.Errorf("RuneReader.EatWhile() error = %v, wantErr %v", gotErr, false)
			}
			if gotCount != tt.wantCount {
				t.Errorf("RuneReader.EatWhile() count = %v, wantCount %v", gotCount, tt.wantCount)
			}

			gotPos := raw.GetPosition()
			if !reflect.DeepEqual(gotPos, tt.wantPos) {
				t.Errorf("RuneReader.GetPosition() gotPos = %v, want %v", gotPos, tt.wantPos)
			}
		})
	}
}

func TestRuneReader_EatWhileMatch(t *testing.T) {
	tests := []struct {
		name      string
		re        *regexp.Regexp
		wantCount int
		wantPos   ReaderPosition
	}{
		{
			"eat nothing",
			regexp.MustCompile("^$"),
			0,
			ReaderPosition{0, 0, false},
		},
		{
			"eat letters",
			regexp.MustCompile("^[aA-zZ]+$"),
			4,
			ReaderPosition{0, 4, false},
		},
		{
			"eat everything",
			regexp.MustCompile(".*"),
			55,
			ReaderPosition{7, 0, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := makeTestReader()

			gotCount, err := raw.EatWhileMatch(tt.re)
			if (err != nil) != false {
				t.Errorf("RuneReader.EatWhileMatch() error = %v, wantErr %v", err, false)
				return
			}
			if gotCount != tt.wantCount {
				t.Errorf("RuneReader.EatWhileMatch() gotCount = %v, want %v", gotCount, tt.wantCount)
			}

			gotPos := raw.GetPosition()
			if !reflect.DeepEqual(gotPos, tt.wantPos) {
				t.Errorf("RuneReader.GetPosition() gotPos = %v, want %v", gotPos, tt.wantPos)
			}
		})
	}
}
