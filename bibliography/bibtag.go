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

/*

# evaluates the content of this BiBTag
# FAILS if this Tag is already evaluated
# returns a list of items which have failed to evaluate
sub evaluate {
  my ($self, %context) = @_;

  my @failed = ();

  # if we have a name, we need to normalize it
  $$self{name}->normalizeValue if defined($$self{name});

  # we need to expand the value and iterate over it
  my @content = @{ $$self{content} };
  return unless scalar(@content) > 0;

  my $item = shift(@content);
  push(@failed, $item->copy) unless $item->evaluate(%context);

  # evaluate and append each content item
  # from the ones that we have
  # DOES NOT DO ANY TYPE CHECKING
  my $cont;
  foreach $cont (@content) {
    push(@failed, $cont) unless $cont->evaluate(%context);
    $item->append($cont);
  }

  # and set the new content
  $$self{content} = $item;

  return @failed;
}
*/
