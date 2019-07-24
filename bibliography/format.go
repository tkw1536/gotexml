package bibliography

import "strings"

// Format represents the format used to format a Bibliography file
type Format struct {
	TagSpace string // spaces around individual parts of Tags

	FirstTagSeperator string // the first seperator
	TagSeparator      string // seperation between every tag
	EntrySuffix       string // a suffix before a closing entry

	EntryKind       EntryKindFormat // how to handle entry kinds
	EntryKindSuffix string

	RemoveEmptyTags  bool // when set to true, remove entry empty tags
	AddTrailingComma bool // when set to true, add a trailing comma to all entries

	FileSeperator string // seperator between different entries in a file
}

// EntryKindFormat represents how to format the kind of an entry
type EntryKindFormat int

// how to format the entrykind field
const (
	EntryKindUntouched EntryKindFormat = iota // leave the kind field as is
	EntryKindUppercase                        // turn it to upper case
	EntryKindLowercase                        // turn it to lower case
)

func (format *Format) tagSpace() string {
	if format == nil {
		return " "
	}
	return format.TagSpace
}

func (format *Format) firstTagSeperator() string {
	if format == nil {
		return ""
	}
	return format.FirstTagSeperator
}

func (format *Format) tagSeparator() string {
	if format == nil {
		return "\n    "
	}
	return format.TagSeparator
}

func (format *Format) entrySuffix() string {
	if format == nil {
		return "\n"
	}
	return format.EntrySuffix
}

func (format *Format) entryKind() EntryKindFormat {
	if format == nil {
		return EntryKindLowercase
	}
	return format.EntryKind
}

func (format *Format) entryKindSuffix() string {
	if format == nil {
		return ""
	}
	return format.EntryKindSuffix
}

func (format *Format) fileSeperator() string {
	if format == nil {
		return "\n\n"
	}
	return format.FileSeperator
}

func (format *Format) removeEmptyTags() bool {
	if format == nil {
		return true
	}
	return format.RemoveEmptyTags
}

// FormatTag formats a Tag using these options
func (format *Format) FormatTag(tag *BibField) {
	space := format.tagSpace()

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

// FormatEntry formats an entry using these options
func (format *Format) FormatEntry(entry *BibEntry) {
	entry.Prefix.Value = ""

	// format the kind
	switch format.entryKind() {
	case EntryKindLowercase:
		entry.Kind.Value = strings.ToLower(entry.Kind.Value)
	case EntryKindUppercase:
		entry.Kind.Value = strings.ToUpper(entry.Kind.Value)
	}

	// set the entry kind suffix
	entry.KindSuffix.Value = format.entryKindSuffix()
	seperator := format.tagSeparator()

	// if we want to remove empty tags, remove them
	if format.removeEmptyTags() {
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
		format.FormatTag(t)
		if i == 0 {
			t.Prefix.Value = format.firstTagSeperator()
		} else {
			t.Prefix.Value = seperator
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
			lastTag.Elements[last].Suffix.Value = format.entrySuffix()
		}
	} else {
		lastTag.Prefix.Value = format.entrySuffix()
	}

}

// FormatFile formats a BibFile entry according to these options
func (format *Format) FormatFile(file *BibFile) {
	// set the seperator
	prefix := format.fileSeperator()
	for i, e := range file.Entries {
		format.FormatEntry(e)
		if i != 0 {
			e.Prefix.Value = prefix
		}
	}
	file.Suffix.Value = "\n" // hard-code end of file
}
