package bibliography

import (
	"io"
	"unicode"

	"github.com/tkw1536/gotexml/utils"
)

// BibEntry is an entry within a BibFile
type BibEntry struct {
	Prefix BibString `json:"prefix"` // Spaces before the BibString

	Kind       *BibString `json:"kind"`       // the type of this BibEntry, a literal succeeding '@'
	KindSuffix *BibString `json:"kindSuffix"` // spaces behind the kind

	Tags []*BibTag `json:"tags"` // tags contained in this BibEntry

	Source utils.ReaderRange `json:"source"` // source of this bibtag
}

// readEntry reads a BibEntry from reader
// Tags end with '}' as a terminating character.
// when err is io.EOF, no beginning tag was found and only Prefix is populated
// else when err is non-nil, it is an instance of utils.ReaderError
func (entry *BibEntry) readEntry(reader *utils.RuneReader) (err error) {
	// skip ahead until we have an '@' preceeded by a space or the beginning of the string
	hasPrevSpace := true
	entry.Prefix.Value, entry.Prefix.Source, err = reader.ReadWhile(func(r rune) bool {
		if hasPrevSpace && r == '@' {
			return false
		}

		hasPrevSpace = unicode.IsSpace(r)
		return true

	})
	if err != nil {
		err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read entry")
		return
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
	if char != '@' {
		err = utils.NewErrorF(reader, "Expected to find an '@' but got %q", char)
		return
	}
	entry.Source.Start = pos

	// read the literal and the appropriate suffix
	entry.Kind = &BibString{}
	entry.KindSuffix, err = entry.Kind.readLiteral(reader)
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
	for true {
		// read the next tag
		t := &BibTag{}
		err = t.readTag(reader)
		if err != nil {
			err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read entry")
			return
		}

		// append the tag to the known list of tags
		entry.Tags = append(entry.Tags, t)

		// if the entry ended
		if t.Suffix.Value == "}" {
			break
		}
	}

	// the last source entry
	entry.Source.End = pos

	// and return the element
	return
}

// Write writes this BibEntry into a writer
func (entry *BibEntry) Write(writer io.Writer) error {
	if err := entry.Prefix.Write(writer); err != nil {
		return err
	}
	if _, err := writer.Write([]byte("@")); err != nil {
		return err
	}
	if err := entry.Kind.Write(writer); err != nil {
		return err
	}
	if err := entry.KindSuffix.Write(writer); err != nil {
		return err
	}
	if _, err := writer.Write([]byte("{")); err != nil {
		return err
	}
	for _, tag := range entry.Tags {
		if err := tag.Write(writer); err != nil {
			return err
		}
	}

	return nil
}
