package bibliography

import (
	"strings"
	"unicode"

	"github.com/tkw1536/gotexml/utils"
)

// BibString represents a string along with a reference to a BibFile
type BibString struct {
	kind  BibStringKind // the type of Bibstring this is
	value string        // the value of this bibstring

	// range of this token inside the source file
	source utils.ReaderRange
}

// BibStringKind represents the specific type of BibStrings that can occur
type BibStringKind string

// kind of BibStrings that occur in the code
const (
	BibStringOther     BibStringKind = ""          // any other kind of bibstring, mostly spaces and un-finished data
	BibStringLiteral   BibStringKind = "LITERAL"   // BibTex Literals (e.g. example)
	BibStringQuote     BibStringKind = "QUOTE"     // BibTex Quotes (e.g. "example")
	BibStringBracket   BibStringKind = "BRACKET"   // anything in BibTex Brackets, e.g. {ExAmPle}
	BibStringEvaluated BibStringKind = "EVALUATED" // anything that has been evaluated (original information has been lost)
)

// Copy makes a copy of this BibString
func (bs *BibString) Copy() *BibString {
	// TODO: Check if this is used
	return &BibString{
		kind:   bs.kind,
		value:  bs.value,
		source: bs.source,
	}
}

// Kind gets the kind of this BibString
func (bs *BibString) Kind() BibStringKind {
	return bs.kind
}

// Value gets the value of this BibString
func (bs *BibString) Value() string {
	return bs.value
}

// Source returns the source of this BibString
func (bs *BibString) Source() utils.ReaderRange {
	return bs.source
}

// NormalizeValue normalizes the value of this BibString
func (bs *BibString) NormalizeValue() {
	bs.value = strings.ToLower(bs.value)
}

// Evaluate evaluates this BibString inside a context
// returns true iff evaluation was successfull
func (bs *BibString) Evaluate(context map[string]string) bool {
	if bs.kind == BibStringLiteral {
		// grab the value of this string
		value, ok := context[strings.ToLower(bs.value)]
		if !ok {
			return false
		}

		// update the kind and value
		bs.kind = BibStringEvaluated
		bs.value = value
	}
	return true
}

// Append appends the value of another bibstring to this bibstring
// This function does not check the legality of such an operation
func (bs *BibString) Append(other *BibString) {
	bs.kind = BibStringEvaluated
	bs.value += other.value
	bs.source.End = other.source.End
}

// readLiteral reads a BibString of kind BibStringLiteral from the input
// Skips and returns spaces after the BibString.
func readLiteral(reader *utils.RuneReader) (lit BibString, space BibString, rr error) {

	lit.source.Start = reader.Position()
	lit.source.End = lit.source.Start

	// read the next character or bail out when an error or EOF occurs
	char, pos, err := reader.Read()
	if err != nil {
		err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read literal")
		return
	}
	if pos.EOF {
		err = utils.NewErrorF(reader, "Unexpected end of input while attempting to read literal")
		return
	}

	// iterate over sequential non-space sequences
	var cache string
	var source utils.ReaderRange
	for isNotSpecialLiteral(char) {
		// add spaces from the previous iteration
		lit.value += cache + string(char)

		// add the next non-space sequence
		cache, _, err = reader.ReadWhile(isNotSpecialSpaceLiteral)
		if err != nil {
			err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read literal")
			return
		}
		lit.value += cache

		// read the next batch of spaces
		cache, source, err = reader.ReadWhile(unicode.IsSpace)
		if err != nil {
			err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read literal")
			return
		}

		// read the next character or bail out
		char, pos, err = reader.Read()
		if err != nil {
			err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read literal")
			return
		}
		if pos.EOF {
			err = utils.NewErrorF(reader, "Unexpected end of input while attempting to read literal")
			return
		}
	}

	// unread the last char
	reader.Unread(char, pos)

	// store remaining data from the literal
	lit.kind = BibStringLiteral
	lit.source.End = source.Start

	// store remaining data from the spacing
	space.kind = BibStringOther
	space.value = cache
	space.source = source

	return
}

// readBrace reads a BibString of kind BibStringBracket from the input
// Must start with "{" and end with "}". Does not skip any spaces before or after.
func readBrace(reader *utils.RuneReader) (brace BibString, err error) {
	char, pos, err := reader.Read()
	if err != nil {
		err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read braces")
		return
	}
	if pos.EOF {
		err = utils.NewErrorF(reader, "Unexpected end of input while attempting to read brace")
		return
	}
	if char != '{' {
		err = utils.NewErrorF(reader, "Expected to find an '{' but got %q", char)
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
		char, pos, err = reader.Read()
		if err != nil {
			err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read braces")
			return
		}
		if pos.EOF {
			err = utils.NewErrorF(reader, "Unexpected end of input while attempting to read braces")
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
	brace.source.End = reader.Position()

	return
}

// readQuote reads a BibString of kind BibStringQuote from the input
// Must start and end with "s. Does not skip any spaces.
func readQuote(reader *utils.RuneReader) (quote BibString, err error) {
	char, pos, err := reader.Read()
	if err != nil {
		err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read quote")
		return
	}
	if pos.EOF {
		err = utils.NewErrorF(reader, "Unexpected end of input while attempting to read quote")
		return
	}
	if char != '"' {
		err = utils.NewErrorF(reader, "Expected to find an '\"' but got %q", char)
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
		char, pos, err = reader.Read()
		if err != nil {
			err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read quote")
			return
		}
		if pos.EOF {
			err = utils.NewErrorF(reader, "Unexpected end of input while attempting to read quote")
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
	quote.source.End = reader.Position()

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
