package bibliography

import (
	"io"

	"github.com/tkw1536/gotexml/utils"
)

// BibFile represents an entire BiBFile
type BibFile struct {
	Entries []BibEntry `json:"entries"` // The entries of this BibFile

	Source utils.ReaderRange `json:"source"` // source range that contains this BibFile
}

// readFile reads a BibFile
func readFile(reader *utils.RuneReader) (file BibFile, err error) {
	// store the original position
	file.Source.Start = reader.Position()
	file.Source.End = file.Source.Start

	// keep reading entries
	var entry BibEntry
	for true {
		entry, err = readEntry(reader)
		if err == io.EOF { // bail out if there are none left
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
