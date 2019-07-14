package bibliography

import (
	"github.com/tkw1536/gotexml/utils"
)

// BibEntry is an entry within a BibFile
type BibEntry struct {
	tp     BibString         // the type of this entry
	tags   []BibTag          // the tags within this BibEntry
	source utils.ReaderRange // the source of this BibEntry
}

// Type returns the type of this BibEntry
func (be *BibEntry) Type() BibString {
	return be.tp
}

// Tags returns the tags of this BibEntry
func (be *BibEntry) Tags() []BibTag {
	return be.tags
}

// Source returns the source of this entry
func (be *BibEntry) Source() utils.ReaderRange {
	return be.source
}

// Evaluate evalutes this BibEntry inside a context
// TODO: Re-implement
func (be *BibEntry) Evaluate(context map[string]string) (failed map[BibString][]*BibString) {
	// normalize the type of this entry
	be.tp.NormalizeValue()

	// evaluate all the entries
	for i := range be.tags {
		errors := be.tags[i].Evaluate(context)
		// store the entries (if needed)
		if len(errors) > 0 {
			// failed[be.tags[i].name] = errors
		}
	}

	return
}
