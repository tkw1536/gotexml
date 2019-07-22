package bibliography

// TODO: Test me
import (
	"io"
	"unicode"

	"github.com/tkw1536/gotexml/utils"
)

// BibTag represents a single (key, value) pairs within a BibTeX entry
// note that value itself may consist of an arbitrary number of BibStrings, being 0 for key-value pairs and a range for concatinated strings
type BibTag struct {
	Prefix   BibString        `json:"prefix"`   // spaces preceeding this BibTag
	Elements []*BibTagElement `json:"elements"` // elements of this BibTag
	Suffix   BibString        `json:"suffix"`   // the suffix (i.e. terminating character) of this string

	Source utils.ReaderRange `json:"source"` // source range that contains this BibTag
}

// BibTagElement represents one element of this BibTag
type BibTagElement struct {
	Value  *BibString     `json:"value"`          // concrete value of this element
	Suffix *BibString     `json:"suffix"`         // space + concatination value of this element
	Role   TagElementRole `json:"role,omitempty"` // the role of this BibTagElement
}

// TagElementRole is the role of a BibTagElement
type TagElementRole string

//kinds of roles that can occur
const (
	NormalElementRole TagElementRole = ""     // no special role
	KeyElementRole    TagElementRole = "key"  // the name of this key
	TermElementRole   TagElementRole = "term" // a term which will be appended
)

// Empty checks if a tag is empty
func (tag *BibTag) Empty() bool {
	return len(tag.Elements) == 0
}

// readTag reads a (potentially empty) BibTag from reader.
// Tags end with a character ',' or '}'. These are contained in suffix.
// when err is no nil, it is an instance of utils.ReaderError.
func (tag *BibTag) readTag(reader *utils.RuneReader) (err error) {
	// read spaces at the beginning
	tag.Prefix.Value, tag.Prefix.Source, err = reader.ReadWhile(unicode.IsSpace)
	if err != nil {
		err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read tag")
		return
	}

	// state we will constantly update
	var r rune
	var s string
	var pos utils.ReaderPosition
	var loc utils.ReaderRange

	// peek the next char
	r, pos, err = reader.Peek()
	if err != nil {
		err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read tag")
	}
	if pos.EOF {
		err = utils.NewErrorF(reader, "Unexpected end of input while attempting to read tag")
		return
	}

	tag.Source.Start = pos

	hadEqualSign := false       // did we have an equal sign yet?
	shouldAppendSuffix := false // should we append further spaces to the suffix instead of resetting it

	// state: modes that are allowed next
	mayStringNext := true
	mayConcatNext := false
	mayEqualNext := false

	var last *BibTagElement

	// iteratively read characters
	for r != ',' && r != '}' {
		switch r {
		case '=':
			if !mayEqualNext {
				err = utils.NewErrorF(reader, "Unexpected \"=\" while attempting to read tag")
				return
			}

			if err = reader.Eat(); err != nil {
				err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read tag")
				return
			}

			// mark the previous element as as a key
			hadEqualSign = true
			last.Role = KeyElementRole

			// add an '=' to the last suffix
			last.Suffix.Value += string(r)
			last.Suffix.Source.End = pos
			shouldAppendSuffix = true

			// after a string, we may no longer have '='s
			mayStringNext = true
			mayConcatNext = false
			mayEqualNext = false
		case '#':
			if !mayConcatNext {
				err = utils.NewErrorF(reader, "Unexpected \"#\" while attempting to read tag")
				return
			}

			if err = reader.Eat(); err != nil {
				err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read tag")
				return
			}

			// mark the previous element as being a term
			last.Role = TermElementRole

			// add an '#' to the last suffix
			last.Suffix.Value += string(r)
			last.Suffix.Source.End = pos
			shouldAppendSuffix = true

			// after a concat, we may only add another string
			mayStringNext = true
			mayConcatNext = false
			mayEqualNext = false
		case '"':
			if !mayStringNext {
				err = utils.NewErrorF(reader, "Unexpected '\"' while attempting to read tag")
				return
			}

			// append a new element
			last = &BibTagElement{
				Value:  &BibString{},
				Suffix: &BibString{},
			}
			tag.Elements = append(tag.Elements, last)

			// read the quote and append it to the content
			err = last.Value.readQuote(reader)
			if err != nil {
				return
			}
			shouldAppendSuffix = false

			// after a quote, either another string is concatinated
			// or everything is terminated
			mayStringNext = false
			mayConcatNext = true
			mayEqualNext = false

		case '{':
			if !mayStringNext {
				err = utils.NewErrorF(reader, "Unexpected '{' while attempting to read tag")
				return
			}

			// append a new element
			last = &BibTagElement{
				Value:  &BibString{},
				Suffix: &BibString{},
			}
			tag.Elements = append(tag.Elements, last)

			// read the brace and append it to the content
			err = last.Value.readBrace(reader)
			if err != nil {
				return
			}
			shouldAppendSuffix = false

			// after a brace, we may have an equal sign unless we already had one before
			mayStringNext = false
			mayConcatNext = false
			mayEqualNext = !hadEqualSign
		default:
			if !mayStringNext {
				err = utils.NewErrorF(reader, "Unexpected start of literal while attempting to read tag")
				return
			}

			// append a new element
			last = &BibTagElement{
				Value: &BibString{},
			}
			tag.Elements = append(tag.Elements, last)

			// read the literal
			last.Suffix, err = last.Value.readLiteral(reader)
			if err != nil {
				return
			}
			shouldAppendSuffix = true

			// after a brace, we may have an equal sign unless we already had one before
			mayStringNext = false
			mayConcatNext = true
			mayEqualNext = !hadEqualSign
		}

		// eat all the spaces
		s, loc, err = reader.ReadWhile(unicode.IsSpace)
		if err != nil {
			err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read tag")
			return
		}

		// append or use s as the space
		if shouldAppendSuffix {
			last.Suffix.AppendRaw(s, loc)
		} else {
			last.Suffix.Value = s
			last.Suffix.Source = loc
		}

		// peek the next char
		r, pos, err = reader.Peek()
		if err != nil {
			err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read tag")
			return
		}
		if pos.EOF {
			err = utils.NewErrorF(reader, "Unexpected end of input while attempting to read tag")
			return
		}
	}

	// end of source
	tag.Source.End = pos

	// eat away the closing tag
	tag.Suffix.Value = string(r)
	tag.Suffix.Source.Start = pos
	tag.Suffix.Source.End = pos
	reader.Eat()

	// and return
	return
}

// Write writes this BibTag into a writer
func (tag *BibTag) Write(writer io.Writer) error {
	if err := tag.Prefix.Write(writer); err != nil {
		return err
	}
	for _, e := range tag.Elements {
		if err := e.Value.Write(writer); err != nil {
			return err
		}
		if err := e.Suffix.Write(writer); err != nil {
			return err
		}
	}
	if err := tag.Suffix.Write(writer); err != nil {
		return err
	}
	return nil
}
