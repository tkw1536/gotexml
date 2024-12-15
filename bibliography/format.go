package bibliography

import (
	"sort"
	"strings"
)

// Formatter formats parts of a [BibFile]
type Formatter struct {
	FieldSpace string // spaces around individual parts of Tags

	FirstFieldSeparator string // the first separators
	FieldSeparator      string // separation between every field
	EntrySuffix         string // a suffix before a closing entry

	EntryKind       EntryKindFormat // how to handle entry kinds
	EntryKindSuffix string

	RemoveEmptyFields bool // when set to true, remove entry empty tags
	AddTrailingComma  bool // when set to true, add a trailing comma to all entries

	FileSeparator string // separator between different entries in a file

	SortEntries bool // if true, sort entries by their key
}

// EntryKindFormat represents how to format the kind of an entry
type EntryKindFormat int

// how to format the entrykind field
const (
	EntryKindUntouched EntryKindFormat = iota // leave the kind field as is
	EntryKindUppercase                        // turn it to upper case
	EntryKindLowercase                        // turn it to lower case
)

// DefaultFormatter is the default formatter.
var DefaultFormatter = Formatter{
	FieldSpace:          " ",
	FirstFieldSeparator: "",
	FieldSeparator:      "\n    ",
	EntrySuffix:         "\n",
	EntryKind:           EntryKindLowercase,
	EntryKindSuffix:     "",
	FileSeparator:       "\n\n",
	RemoveEmptyFields:   true,
	SortEntries:         true,
}

// Format formats a file according to the options set.
func (format Formatter) Format(file *BibFile) {
	prefix := format.FileSeparator
	for i, e := range file.Entries {
		format.entry(e)
		if i != 0 {
			e.Prefix.Value = prefix
		}
	}

	// sort if requested
	if format.SortEntries {
		format.sort(file)
	}

	file.Suffix.Value = "\n" // hard-code end of file
}

// entry formats the given entry
func (format Formatter) entry(entry *BibEntry) {
	entry.Prefix.Value = ""

	// format the kind
	switch format.EntryKind {
	case EntryKindLowercase:
		entry.Kind.Value = strings.ToLower(entry.Kind.Value)
	case EntryKindUppercase:
		entry.Kind.Value = strings.ToUpper(entry.Kind.Value)
	}

	// set the entry kind suffix
	entry.KindSuffix.Value = format.EntryKindSuffix
	TagSeparator := format.FieldSeparator

	// if we want to remove empty tags, remove them
	if format.RemoveEmptyFields {
		filteredTags := entry.Fields[:0]
		for _, t := range entry.Fields {
			if !t.Empty() {
				filteredTags = append(filteredTags, t)
			}
		}
		for i := len(filteredTags); i < len(entry.Fields); i++ {
			entry.Fields[i] = nil
		}
		entry.Fields = filteredTags
	}

	// format the tags
	for i, t := range entry.Fields {
		format.field(t)
		if i == 0 {
			t.Prefix.Value = format.FirstFieldSeparator
		} else {
			t.Prefix.Value = TagSeparator
		}
	}

	// we now need to format the last tag in the entry
	// but this can't be done if there are no tags
	last := len(entry.Fields) - 1
	if last == -1 {
		return
	}
	lastTag := entry.Fields[last]

	// make sure it ends with a '}' (if we filtered)
	lastTag.Suffix.Value = "}"

	// set it as the suffix of the last element
	// if the tag is empty
	if !lastTag.Empty() {
		last = len(lastTag.Elements) - 1
		if last != -1 {
			lastTag.Elements[last].Suffix.Value = format.EntrySuffix
		}
	} else {
		lastTag.Prefix.Value = format.EntrySuffix
	}

}

// field formats a single field
func (format Formatter) field(tag *BibField) {
	space := format.FieldSpace

	// Prefix and suffix set by FormatEntry()
	tag.Prefix.Value = ""

	for _, e := range tag.Elements {
		switch e.Role {
		case NormalElementRole:
			e.Suffix.Value = ""
		case KeyElementRole:
			e.Suffix.Value = space + "=" + space
		case TermElementRole:
			e.Suffix.Value = space + "#" + space
		}
	}
}

// sort sorts the entries in the given file.
// This will sort even if [SortEntries] is unset.
func (format Formatter) sort(file *BibFile) {
	// get the keys for each entry
	keys := make(map[*BibEntry]string, len(file.Entries))
	for _, entry := range file.Entries {
		keys[entry] = entry.Label()
	}

	// and sort by the keys!
	sort.Slice(file.Entries, func(i, j int) bool {
		return keys[file.Entries[i]] < keys[file.Entries[j]]
	})
}
