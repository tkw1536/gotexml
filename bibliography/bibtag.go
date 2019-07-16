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
	Prefix   BibString       `json:"prefix"`   // spaces preceeding this BibTag
	Elements []BibTagElement `json:"elements"` // elements of this BibTag

	Source utils.ReaderRange `json:"source"` // source range that contains this BibTag
}

// BibTagElement represents one element of this BibTag
type BibTagElement struct {
	Value        BibString // concrete value of this element
	Suffix       BibString // space + concatination value of this element
	IsKeyElement bool      `json:"isKeyElement,omitempty"` // true iff this is a key element
}

// reads a BibTag (which does not include a terminating character) from the source file
// contains all spaces before a terminating "," or "}"
// err is either nil, io.EOF or an instance of
// when io.EOF is returned, this means that no valid BibTag was read and only tag.initialSpace was populated
func readTag(reader *utils.RuneReader) (tag BibTag, err error) {
	// read spaces at the beginning
	tag.Prefix.Value, tag.Prefix.Source, err = reader.ReadWhile(unicode.IsSpace)
	if err != nil {
		err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read tag")
		return
	}

	// state we will constantly update
	var r rune
	var s string
	var temp, space BibString
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

	// check if we need to exit
	if r == '}' || r == ',' {
		err = io.EOF
		return
	}
	tag.Source.Start = pos
	prevPos := pos

	hadEqualSign := false       // did we have an equal sign yet?
	shouldAppendSuffix := false // should we append further spaces to the suffix instead of resetting it

	// state: modes that are allowed next
	mayStringNext := true
	mayConcatNext := false
	mayEqualNext := false

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

			// we had an equal sign, so we may now set '='
			hadEqualSign = true
			tag.Elements[0].IsKeyElement = true

			// add an '=' to the last suffix
			tag.Elements[len(tag.Elements)-1].Suffix.Value += string(r)
			tag.Elements[len(tag.Elements)-1].Suffix.Source.End = pos
			prevPos = pos
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

			// add an '#' to the last suffix
			tag.Elements[len(tag.Elements)-1].Suffix.Value += string(r)
			tag.Elements[len(tag.Elements)-1].Suffix.Source.End = pos
			prevPos = pos
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

			// read the quote and append it to the content
			temp, err = readQuote(reader)
			if err != nil {
				return
			}
			tag.Elements = append(tag.Elements, BibTagElement{
				Value: temp,
			})
			prevPos = temp.Source.End
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

			// read the brace and append it to the content
			temp, err = readBrace(reader)
			if err != nil {
				return
			}
			tag.Elements = append(tag.Elements, BibTagElement{
				Value: temp,
			})
			prevPos = temp.Source.End
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

			// read the literal
			temp, space, err = readLiteral(reader)
			if err != nil {
				return
			}
			tag.Elements = append(tag.Elements, BibTagElement{
				Value:  temp,
				Suffix: space,
			})

			prevPos = temp.Source.End
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
		l := &tag.Elements[len(tag.Elements)-1].Suffix
		if shouldAppendSuffix {
			l.AppendRaw(s, loc)
		} else {
			l.Value = s
			l.Source = loc
		}
		tag.Elements[len(tag.Elements)-1].Suffix = *l

		// if we read a space, update the last read position
		if s != "" {
			prevPos = l.Source.End
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

	tag.Source.End = prevPos
	return
}
