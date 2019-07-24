package bibliography

// TODO: Test me
import (
	"io"
	"unicode"

	"github.com/tkw1536/gotexml/utils"
)

// BibField represents a single (key, value) pairs within a BibTeX entry
// note that value itself may consist of an arbitrary number of BibStrings, being 0 for key-value pairs and a range for concatinated strings
type BibField struct {
	Prefix   BibString          `json:"prefix"`   // spaces preceeding this BibField
	Elements []*BibFieldElement `json:"elements"` // elements of this BibField
	Suffix   BibString          `json:"suffix"`   // the suffix (i.e. terminating character) of this string

	Source utils.ReaderRange `json:"source"` // source range that contains this BibField
}

// BibFieldElement represents one element of this BibField
type BibFieldElement struct {
	Value  *BibString       `json:"value"`          // concrete value of this element
	Suffix *BibString       `json:"suffix"`         // space + concatination value of this element
	Role   FieldElementRole `json:"role,omitempty"` // the role of this BibFieldElement
}

// FieldElementRole is the role of a BibFieldElement
type FieldElementRole string

//kinds of roles that can occur
const (
	NormalElementRole FieldElementRole = ""     // no special role
	KeyElementRole    FieldElementRole = "key"  // the name of this key
	TermElementRole   FieldElementRole = "term" // a term which will be appended
)

// readField reads a (potentially empty) BibField from reader.
// Fields end with a character ',' or '}'. These are contained in suffix.
// when err is no nil, it is an instance of utils.ReaderError.
func (field *BibField) readField(reader *utils.RuneReader) (err error) {
	// read spaces at the beginning
	field.Prefix.Value, field.Prefix.Source, err = reader.ReadWhile(unicode.IsSpace)
	if err != nil {
		err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read field")
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
		err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read field")
	}
	if pos.EOF {
		err = utils.NewErrorF(reader, "Unexpected end of input while attempting to read field")
		return
	}

	field.Source.Start = pos

	hadEqualSign := false       // did we have an equal sign yet?
	shouldAppendSuffix := false // should we append further spaces to the suffix instead of resetting it

	// state: modes that are allowed next
	mayStringNext := true
	mayConcatNext := false
	mayEqualNext := false

	var last *BibFieldElement

	// iteratively read characters
	for r != ',' && r != '}' {
		switch r {
		case '=':
			if !mayEqualNext {
				err = utils.NewErrorF(reader, "Unexpected \"=\" while attempting to read field")
				return
			}

			if err = reader.Eat(); err != nil {
				err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read field")
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
				err = utils.NewErrorF(reader, "Unexpected \"#\" while attempting to read field")
				return
			}

			if err = reader.Eat(); err != nil {
				err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read field")
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
				err = utils.NewErrorF(reader, "Unexpected '\"' while attempting to read field")
				return
			}

			// append a new element
			last = &BibFieldElement{
				Value:  &BibString{},
				Suffix: &BibString{},
			}
			field.Elements = append(field.Elements, last)

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
				err = utils.NewErrorF(reader, "Unexpected '{' while attempting to read field")
				return
			}

			// append a new element
			last = &BibFieldElement{
				Value:  &BibString{},
				Suffix: &BibString{},
			}
			field.Elements = append(field.Elements, last)

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
				err = utils.NewErrorF(reader, "Unexpected start of literal while attempting to read field")
				return
			}

			// append a new element
			last = &BibFieldElement{
				Value: &BibString{},
			}
			field.Elements = append(field.Elements, last)

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
			err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read field")
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
			err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read field")
			return
		}
		if pos.EOF {
			err = utils.NewErrorF(reader, "Unexpected end of input while attempting to read field")
			return
		}
	}

	// end of source
	field.Source.End = pos

	// eat away the closing field
	field.Suffix.Value = string(r)
	field.Suffix.Source.Start = pos
	field.Suffix.Source.End = pos
	reader.Eat()

	// and return
	return
}

// Write writes this BibField into a writer
func (field *BibField) Write(writer io.Writer) error {
	if err := field.Prefix.Write(writer); err != nil {
		return err
	}
	for _, e := range field.Elements {
		if err := e.Value.Write(writer); err != nil {
			return err
		}
		if err := e.Suffix.Write(writer); err != nil {
			return err
		}
	}
	if err := field.Suffix.Write(writer); err != nil {
		return err
	}
	return nil
}

// Empty checks if a field is empty
func (field *BibField) Empty() bool {
	return len(field.Elements) == 0
}

// IsKeyValue checks if this BibField is of the form 'key = value'
func (field *BibField) IsKeyValue() bool {
	return len(field.Elements) >= 1 && field.Elements[0].Role == KeyElementRole
}

// GetKey returns the 'key' of this BibField entry, i.e. the first element in a 'key = value' assignment
// if there is no key, returns nil
func (field *BibField) GetKey() *BibFieldElement {
	if !field.IsKeyValue() {
		return nil
	}

	return field.Elements[0]
}

// GetValue returns the value elements of this key, i.e. everything after the first element in a 'key = value' assignment
// if the BibEntry is not of the form key == value, returns 0
func (field *BibField) GetValue() []*BibFieldElement {
	if !field.IsKeyValue() {
		return nil
	}

	return field.Elements[1:]
}
