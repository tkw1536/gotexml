package bibliography

import (
	"io"
	"unicode"

	"github.com/tkw1536/gotexml/utils"
)

// BibEntry is an entry within a BibFile
type BibEntry struct {
	prefix BibString // spaces before everything

	kind       BibString // the type (e.g. 'string' for '@string{}')
	kindSuffix BibString // the (space) suffix for this entry

	tags []BibTag // elements of this BibTag

	source utils.ReaderRange // source of this bibtag
}

// Source returns the source of this BibEntry
func (be *BibEntry) Source() utils.ReaderRange {
	return be.source
}

func readEntry(reader *utils.RuneReader) (entry BibEntry, err error) {
	// skip ahead until we have an '@' preceeded by a space or the beginning of the string
	hasPrevSpace := true
	entry.prefix.value, entry.prefix.source, err = reader.ReadWhile(func(r rune) (c bool) {
		c = !hasPrevSpace || r != '@'
		hasPrevSpace = unicode.IsSpace(r)
		return
	})
	if err != nil {
		err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read entry")
	}

	// read an '@' sign or bail out
	char, pos, err := reader.Read()
	if err != nil {
		err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read entry")
		return
	}
	if pos.EOF {
		err = io.EOF
		return
	}
	if char != '{' {
		err = utils.NewErrorF(reader, "Expected to find an '@' but got %q", char)
		return
	}
	entry.source.Start = pos

	// read the literal and the appropriate suffix
	entry.kind, entry.kindSuffix, err = readLiteral(reader)
	if err != nil {
		err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read entry")
		return
	}

	// read a '{' or bail out
	char, pos, err = reader.Read()
	if err != nil {
		err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read entry")
		return
	}
	if pos.EOF {
		err = utils.NewErrorF(reader, "Unexpected end of input while attempting to read entry")
		return
	}
	if char != '{' {
		err = utils.NewErrorF(reader, "Expected to find an '{' but got %q", char)
		return
	}

	// continously read tags from this entry
	// until we have an io.EOF error reported
	var t BibTag
	for true {
		// read the next tag and break on EOF
		t, err = readTag(reader)
		if err == io.EOF {
			break
		}
		if err != nil {
			err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read entry")
			return
		}

		// append the tag to the known list of tags
		entry.tags = append(entry.tags, t)

		// read the next char
		char, pos, err = reader.Read()
		if err != nil {
			err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read entry")
			return
		}
		if pos.EOF {
			err = utils.NewErrorF(reader, "Unexpected end of input while attempting to read entry")
			return
		}

		// } => end of entry
		if char == '}' {
			break
		}

		// if we don't have a "," (i.e. next entry), something went wrong
		if char != ',' {
			err = utils.NewErrorF(reader, "Expected to find a ',' or '}' but got %q", char)
		}
	}

	// the last source entry
	entry.source.End = pos

	// and return the element
	return
}
