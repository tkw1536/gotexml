package bibliography

// TODO: Test me
import (
	"io"
	"unicode"

	"github.com/tkw1536/gotexml/utils"
)

// BibTag represents a tag (e.g. author = {1 2 3}) inside of a BibFile
type BibTag struct {
	prefix   BibString       // spaces before everything
	elements []BibTagElement // elements of this BibTag

	source utils.ReaderRange // source of this bibtag
}

// BibTagElement represents one element of this BibTag
type BibTagElement struct {
	isKeyElement bool

	name   BibString
	suffix BibString
}

// Prefix returns the name of this BibTag
func (bt *BibTag) Prefix() BibString {
	return bt.prefix
}

// Elements returns the elements of this BibTag
func (bt *BibTag) Elements() []BibTagElement {
	return bt.elements
}

// Source returns the source of this BibTag
func (bt *BibTag) Source() utils.ReaderRange {
	return bt.source
}

// reads a BibTag (which does not include a terminating character) from the source file
// contains all spaces before a terminating "," or "}"
// err is either nil, io.EOF or an instance of
// when io.EOF is returned, this means that no valid BibTag was read and only tag.initialSpace was populated
func readTag(reader *utils.RuneReader) (tag BibTag, err error) {
	// read spaces at the beginning
	tag.prefix.value, tag.prefix.source, err = reader.ReadWhile(unicode.IsSpace)
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
	tag.source.Start = pos
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
			tag.elements[0].isKeyElement = true

			// add an '=' to the last suffix
			tag.elements[len(tag.elements)-1].suffix.value += string(r)
			tag.elements[len(tag.elements)-1].suffix.source.End = pos
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
			tag.elements[len(tag.elements)-1].suffix.value += string(r)
			tag.elements[len(tag.elements)-1].suffix.source.End = pos
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
			tag.elements = append(tag.elements, BibTagElement{
				name: temp,
			})
			prevPos = temp.source.End
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
			tag.elements = append(tag.elements, BibTagElement{
				name: temp,
			})
			prevPos = temp.source.End
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
			tag.elements = append(tag.elements, BibTagElement{
				name:   temp,
				suffix: space,
			})

			prevPos = temp.source.End
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
		l := &tag.elements[len(tag.elements)-1].suffix
		if shouldAppendSuffix {
			l.AppendRaw(s, loc)
		} else {
			l.value = s
			l.source = loc
		}
		tag.elements[len(tag.elements)-1].suffix = *l

		// if we read a space, update the last read position
		if s != "" {
			prevPos = l.source.End
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

	tag.source.End = prevPos
	return
}
