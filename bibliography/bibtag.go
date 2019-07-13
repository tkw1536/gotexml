package bibliography

// TODO: Test me
import (
	"unicode"

	"github.com/tkw1536/gotexml/utils"
)

// BibTag represents a tag (e.g. author = {1 2 3}) inside of a BibFile
type BibTag struct {
	initialSpace BibString // spaces before everything

	name    BibString   // the name of this tag
	content []BibString // the content of this tag

	space  BibString
	source utils.ReaderRange // source of this bibtag
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

// reads a single tag from the reader, with optional name and content
// skips only spaces before the input but not after
func readTag(reader *utils.RuneReader) (tag BibTag, err error) {
	// read spaces at the beginning
	tag.initialSpace.value, tag.initialSpace.source, err = reader.ReadWhile(unicode.IsSpace)

	//
	var r rune
	r, tag.source.Start, err = reader.Read()
	if err != nil {
		err = utils.WrapErrorF(reader, err, "Unable to read")
	}

	/*
		my ($reader) = @_;

		# skip spaces and start reading a tag
		$reader->eatSpaces;
		my ($sr, $sc) = $reader->getPosition;
		my ($er, $ec) = ($sr, $sc);

		# if we only have a closing brace
		# we may have tried to read a closing brace
		# so return undef and also no error.
		my ($char) = $reader->peekChar;
		return undef, 'unexpected end of input while reading tag', $reader->getLocation unless defined($char);

		if ($char eq '}' or $char eq ',') {
		  return undef, undef;
		}

		# STATE: What we are allowed to read next
		my $mayStringNext = 1;
		my $mayConcatNext = 0;
		my $mayEqualNext  = 0;

		# results and if we had an error
		my @content = ();
		my ($value, $valueError, $valueLocation);
		my $hadEqualSign = 0;

		# read until we encounter a , or a closing brace
		while ($char ne ',' && $char ne '}') {

		  # if we have an equals sign, remember that we had one
		  # and allow only strings next (i.e. the value)
		  if ($char eq '=') {
			return undef, 'unexpected "="', $reader->getLocation unless $mayEqualNext;
			$reader->eatChar;

			$hadEqualSign = 1;

			$mayStringNext = 1;
			$mayConcatNext = 0;
			$mayEqualNext  = 0;

			# if we have a concat, allow only strings (i.e. the value) next
		  } elsif ($char eq '#') {
			return undef, 'unexpected "#"', $reader->getLocation unless $mayConcatNext;
			$reader->eatChar;

			$mayStringNext = 1;
			$mayConcatNext = 0;
			$mayEqualNext  = 0;

			# if we had a quote, allow only a concat next
		  } elsif ($char eq '"') {
			return undef, 'unexpected \'"\'', $reader->getLocation unless $mayStringNext;

			($value, $valueError, $valueLocation) = readQuote($reader);
			return $value, $valueError, $valueLocation unless defined($value);
			push(@content, $value);

			$mayStringNext = 0;
			$mayConcatNext = 1;
			$mayEqualNext  = 0;

			# if we had a brace, allow only a concat next
		  } elsif ($char eq '{') {
			return undef, 'unexpected \'{\'', $reader->getLocation unless $mayStringNext;

			($value, $valueError, $valueLocation) = readBrace($reader);
			return $value, $valueError, $valueLocation unless defined($value);
			push(@content, $value);

			$mayStringNext = 0;
			$mayConcatNext = 0;
			$mayEqualNext  = !$hadEqualSign;

			# if we have a literal, allow concat and equals next (unless we already had)
		  } else {
			return undef, 'unexpected start of literal', $reader->getPosition unless $mayStringNext;

			($value, $valueError, $valueLocation) = readLiteral($reader);
			return $value, $valueError, $valueLocation unless defined($value);
			push(@content, $value);

			$mayStringNext = 0;
			$mayConcatNext = 1;
			$mayEqualNext  = !$hadEqualSign;
		  }

		  ($er, $ec) = $reader->getPosition;
		  $reader->eatSpaces;

		  ($char) = $reader->peekChar;
		  return undef, 'unexpected end of input while reading tag', $reader->getPosition unless defined($char);
		}

		# if we had an equal sign, shift that value
		my $name;
		$name = shift(@content) if ($hadEqualSign);

		return BiBTeXML::Bibliography::BibTag->new($name, [@content], [($sr, $sc, $er, $ec)]);
	*/
	return
}
