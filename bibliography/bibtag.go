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

// Value returns the singular value of this BibString or panic()s
func (bt *BibTag) Value() BibString {
	if len(bt.elements) != 2 {
		panic("BibTag.Value() is not evaluated")
	}
	return bt.elements[1].name
}

// Evaluate evaluates the content of this Bibtag
// returns a list of items which have failed to evaluate
func (bt *BibTag) Evaluate(context map[string]string) (failed []*BibString) {
	/*
		// normalize the value of the name
		bt.name.NormalizeValue()

		// if we do not have any content, we panic(
		if len(bt.content) == 0 {
			panic("BibTag.Evaluate() is of length 0")
		}

		// evaluate the first element
		// and append it to failed if that fails
		evaluated := bt.content[0]
		if !evaluated.Evaluate(context) {
			failed = append(failed, evaluated.Copy())
		}

		// evaluate and then append all subsequent items
		for _, item := range bt.content[1:] {
			if !item.Evaluate(context) {
				failed = append(failed, &item)
			}
			evaluated.Append(&item)
		}

		// store the evaluated content
		bt.content = []BibString{evaluated}
		return
	*/
	return
}

// reads a BibTag and a suffix (terminating character, e.g. '}', from the source file
// err is either nil, io.EOF or an instance of
// when io.EOF is returned, this means that no valid BibTag was read and only tag.initialSpace was populated
func readTag(reader *utils.RuneReader) (tag BibTag, err error) {
	// TODO: Spaces do not seem accurate for suffixes
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
		}

		// append the new spaces to the suffix if we need to
		if shouldAppendSuffix {
			tag.elements[len(tag.elements)-1].suffix.value += s
			tag.elements[len(tag.elements)-1].suffix.source.End = loc.End
		} else {
			tag.elements[len(tag.elements)-1].suffix.value = s
			tag.elements[len(tag.elements)-1].suffix.source = loc
		}

		// peek the next char
		r, pos, err = reader.Peek()
		if err != nil {
			err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read tag")
		}
		if pos.EOF {
			err = utils.NewErrorF(reader, "Unexpected end of input while attempting to read tag")
			return
		}
	}

	tag.source.End = reader.Position()
	return
}
