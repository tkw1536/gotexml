package bibliography

import (
	"io"

	"github.com/tkw1536/gotexml/utils"
)

// BibFile represents an entire BiBFile
type BibFile struct {
	Entries []*BibEntry `json:"entries"` // The entries of this BibFile
	Suffix  BibString   `json:"suffix"`  // the suffix of this file

	Source utils.ReaderRange `json:"source"` // source range that contains this BibFile
}

// NewBibFileFromReader makes a new BibFile from the given reader
func NewBibFileFromReader(reader *utils.RuneReader) (file *BibFile, err error) {
	file = &BibFile{}
	err = file.readFile(reader)
	return
}

// readFile reads a BibFile from reader
func (file *BibFile) readFile(reader *utils.RuneReader) (err error) {
	// store the original position
	file.Source.Start = reader.Position()
	file.Source.End = file.Source.Start

	// keep reading entries

	for true {
		entry := &BibEntry{}
		err = entry.readEntry(reader)
		if err == io.EOF { // bail out if there are none left
			file.Suffix = entry.Prefix
			err = nil
			break
		}

		// throw an error
		if err != nil {
			err = utils.WrapErrorF(reader, err, "Unexpected error while attempting to read tag")
			return
		}

		file.Entries = append(file.Entries, entry)
		file.Source.End = entry.Source.End
	}
	return
}

// Write writes this BibFile into a writer
func (file *BibFile) Write(writer io.Writer) error {
	for _, e := range file.Entries {
		if err := e.Write(writer); err != nil {
			return err
		}
	}
	return file.Suffix.Write(writer)
}
