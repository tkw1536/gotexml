package bibliography

import (
	"strings"
	"unicode"

	"github.com/pkg/errors"
	"github.com/tkw1536/gotexml/utils"
)

// BibFile represents an entire BibFile
type BibFile struct {
	Context map[string]string
	Reader  *utils.RuneReader

	entries []BibEntry // parsed entries
}

// Entries returns the entries of this BiBFile
func (bf *BibFile) Entries() []BibEntry {
	return bf.entries
}

// Parse Parses a BiBFile
func (bf *BibFile) Parse() error {
	return nil
}

// readLiteral reads a BibString of kind BibStringLiteral from the input
// skips spaces after the read literal
func (bf *BibFile) readLiteral() (lit BibString, err error) {

	lit.source.Start = bf.Reader.GetPosition()
	lit.source.End = bf.Reader.GetPosition()

	// read the next character or bail out when an error or EOF occurs
	char, pos, err := bf.Reader.Read()
	if err != nil {
		err = errors.Wrapf(err, "Unexpected error while attempting to read literal near %s", bf.Reader.GetPosition())
		return
	}
	if pos.EOF {
		err = errors.Errorf("Unexpected end of input while attempting to read literal near %s", pos)
		return
	}

	// iterate over sequential non-space sequences
	var cache string
	var source utils.ReaderRange
	for isNotSpecialLiteral(char) {
		// add spaces from the previous iteration
		lit.value += cache + string(char)

		// add the next non-space sequence
		cache, _, err = bf.Reader.ReadWhile(isNotSpecialSpaceLiteral)
		if err != nil {
			err = errors.Wrapf(err, "Unexpected error while attempting to read literal near %s", bf.Reader.GetPosition())
			return
		}
		lit.value += cache

		// read the next batch of spaces
		cache, source, err = bf.Reader.ReadWhile(unicode.IsSpace)
		lit.source.End = source.Start
		if err != nil {
			err = errors.Wrapf(err, "Unexpected error while attempting to read literal near %s", bf.Reader.GetPosition())
			return
		}

		// read the next character or bail out
		char, pos, err = bf.Reader.Read()
		if err != nil {
			err = errors.Wrapf(err, "Unexpected error while attempting to read literal near %s", bf.Reader.GetPosition())
			return
		}
		if pos.EOF {
			err = errors.Errorf("Unexpected end of input while attempting to read literal near %s", pos)
			return
		}
	}

	// unread the last char
	bf.Reader.Unread(char, pos)

	// and return the literal
	lit.kind = BibStringLiteral
	return
}

// readBrace reads a BibString of kind BibStringBracket from the input
// does not skip any spaces before or after
func (bf *BibFile) readBrace() (brace BibString, err error) {
	char, pos, err := bf.Reader.Read()
	if err != nil {
		err = errors.Wrapf(err, "Unexpected error while attempting to read braces near %s", bf.Reader.GetPosition())
		return
	}
	if pos.EOF {
		err = errors.Errorf("Unexpected end of input while attempting to read brace near %s", pos)
		return
	}
	if char != '{' {
		err = errors.Errorf("Expected to find an '{' near %s but got %q", pos, char)
		return
	}

	// record starting position
	brace.source.Start = pos

	// iteratively read chars, keeping track of the current level
	var builder strings.Builder
	level := 1
	for true {
		// read the next character
		// and bail out when an error or EOF occurs
		char, pos, err = bf.Reader.Read()
		if err != nil {
			err = errors.Wrapf(err, "Unexpected error while attempting to read braces near %s", bf.Reader.GetPosition())
			return
		}
		if pos.EOF {
			err = errors.Errorf("Unexpected end of input while attempting to read braces near %s", pos)
			return
		}

		// update level
		if char == '{' {
			level++
		} else if char == '}' {
			level--
		}

		// final closing brace => exit
		if level == 0 {
			break
		}

		// record the rune
		builder.WriteRune(char)

	}

	brace.kind = BibStringBracket
	brace.value = builder.String()
	brace.source.End = bf.Reader.GetPosition()

	return
}

// readQuote reads a BibString of kind BibStringQuote from the input
// does not skip any spaces before or after
func (bf *BibFile) readQuote() (quote BibString, err error) {
	char, pos, err := bf.Reader.Read()
	if err != nil {
		err = errors.Wrapf(err, "Unexpected error while attempting to read quote near %s", bf.Reader.GetPosition())
		return
	}
	if pos.EOF {
		err = errors.Errorf("Unexpected end of input while attempting to read quote near %s", pos)
		return
	}
	if char != '"' {
		err = errors.Errorf("Expected to find an '\"' near %s but got %q", pos, char)
		return
	}

	// record starting position
	quote.source.Start = pos

	// iteratively read chars, keeping track of the current level
	var builder strings.Builder
	level := 0
	for true {
		// read the next character
		// and bail out when an error or EOF occurs
		char, pos, err = bf.Reader.Read()
		if err != nil {
			err = errors.Wrapf(err, "Unexpected error while attempting to read quote near %s", bf.Reader.GetPosition())
			return
		}
		if pos.EOF {
			err = errors.Errorf("Unexpected end of input while attempting to read quote near %s", pos)
			return
		}

		if char == '"' {
			if level == 0 {
				break
			}
		} else if char == '{' {
			level++
		} else if char == '}' {
			// if we are at level 0, we ignore extra closing braces
			// as being an error
			if level > 0 {
				level--
			}
		}

		builder.WriteRune(char)
	}

	quote.kind = BibStringQuote
	quote.value = builder.String()
	quote.source.End = bf.Reader.GetPosition()

	return
}

// isNotSpecialLiteral checks that the rune is not a specially treated literal
func isNotSpecialLiteral(r rune) bool {
	return (r != '{') && (r != '}') && (r != '=') && (r != '#') && (r != ',')
}

// isNotSpecialSpaceLiteral checks the the rune is not a specially treated literal
// and also does not terminate a space
func isNotSpecialSpaceLiteral(r rune) bool {
	return !unicode.IsSpace(r) && (r != '{') && (r != '}') && (r != '=') && (r != '#') && (r != ',')
}
