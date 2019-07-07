package bibliography

// TODO: Test me
import "github.com/tkw1536/gotexml/utils"

// BibTag represents a tag (e.g. author = {1 2 3}) inside of a BibFile
type BibTag struct {
	name    BibString         // the name of this tag
	content []BibString       // the content of this tag
	source  utils.ReaderRange // source of this bibtag
}

// Name returns the name of this BibTag
func (bt *BibTag) Name() BibString {
	return bt.name
}

// Content returns the content of this BibTag
func (bt *BibTag) Content() []BibString {
	return bt.content
}

// Source returns the source of this BibTag
func (bt *BibTag) Source() utils.ReaderRange {
	return bt.source
}

// Value returns the singular value of this BibString or panic()s
func (bt *BibTag) Value() BibString {
	if len(bt.content) != 1 {
		panic("BibTag.Value() is not evaluated")
	}
	return bt.content[0]
}

// Evaluate evaluates the content of this Bibtag
// returns a list of items which have failed to evaluate
func (bt *BibTag) Evaluate(context map[string]string) (failed []*BibString) {
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
}
