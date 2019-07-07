package bibliography

import (
	"strings"

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
	BibStringOther     BibStringKind = ""
	BibStringLiteral   BibStringKind = "LITERAL"
	BibStringQuote     BibStringKind = "QUOTE"
	BibStringBracket   BibStringKind = "BRACKET"
	BibStringEvaluated BibStringKind = "EVALUATED"
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
