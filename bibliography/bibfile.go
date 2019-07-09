package bibliography

import (
	"github.com/pkg/errors"
	"github.com/tkw1536/gotexml/utils"
)

// BibFile represents an entire BibFile
type BibFile struct {
	Context map[string]string
	Reader  *utils.RuneReader

	entries []BibEntry // parsed entries
}

// Entries returns the entries of this BiBFile
func (bf *BibFile) Entries() []BibEntry {
	return bf.entries
}

// Parse Parses a BiBFile
func (bf *BibFile) Parse() error {
	return errors.New("Not implemented")
}
