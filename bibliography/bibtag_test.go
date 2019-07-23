package bibliography

import (
	"bytes"
	"path"
	"reflect"
	"testing"

	"github.com/tkw1536/gotexml/utils"
)

func Test_readTag(t *testing.T) {
	tests := []struct {
		name  string
		input string
		asset string
	}{
		// value only
		{"empty tag", ``, "0001_empty"},
		{"endingbrace", `}`, "0010_endingbrace"},
		{"literal value", `value`, "0002_literal"},
		{"quoted value", `"value"`, "0003_quoted"},
		{"braced value", `{value}`, "0004_braced"},
		{"concated literals", `value1 # value2`, "0005_concated"},
		{"concated quote and literal", `"value1" # value2`, "0006_concat_quote_literal"},
		// key = value
		{"simple name", `name = value`, "0007_key_value"},
		{"simple name (compact)", `name=value`, "0008_simple_name_compact"},
		{"name + compact value", `name=a#"b"`, "0009_name_compact_value"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// read the assets
			var wantTag *BibTag
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibtag_read", tt.asset+".json"), &wantTag)

			// call readTag
			gotTag := &BibTag{}
			err := gotTag.readTag(utils.NewRuneReaderFromString(tt.input + ", "))

			if (err != nil) != false {
				t.Errorf("BibTag.readTag() error = %v, wantErr %v", err, false)
				return
			}

			if !reflect.DeepEqual(gotTag, wantTag) {
				t.Errorf("BibTag.readTag() = %v, want %v", gotTag, wantTag)
			}
		})
	}
}

func Benchmark_ReadTag_Empty(b *testing.B) {
	benchmarkReadTag(``, b)
}

func Benchmark_ReadTag_Value(b *testing.B) {
	benchmarkReadTag(`value`, b)
}

func Benchmark_ReadTag_QuotedValue(b *testing.B) {
	benchmarkReadTag(`"value"`, b)
}

func Benchmark_ReadTag_BracketedValue(b *testing.B) {
	benchmarkReadTag(`{value}`, b)
}

func Benchmark_ReadTag_Concat(b *testing.B) {
	benchmarkReadTag(`value1 # value2`, b)
}

func Benchmark_ReadTag_ComplexConcat(b *testing.B) {
	benchmarkReadTag(`"value1" # value2`, b)
}

func Benchmark_ReadTag_KeyValue(b *testing.B) {
	benchmarkReadTag(`name = value`, b)
}

func Benchmark_ReadTag_KeyValueCompact(b *testing.B) {
	benchmarkReadTag(`name=value`, b)
}

func Benchmark_ReadTag_KeyValueComplex(b *testing.B) {
	benchmarkReadTag(`name=a#"b"`, b)
}

func benchmarkReadTag(content string, b *testing.B) {
	p := content + ", "
	tag := &BibTag{}
	for n := 0; n < b.N; n++ {
		tag.readTag(utils.NewRuneReaderFromString(p))
	}
}

func TestBibTag_Write(t *testing.T) {
	tests := []struct {
		name       string
		wantString string
		asset      string
	}{
		// value only
		{"empty tag", `,`, "0001_empty"},
		{"endingbrace", `}`, "0010_endingbrace"},
		{"literal value", `value,`, "0002_literal"},
		{"quoted value", `"value",`, "0003_quoted"},
		{"braced value", `{value},`, "0004_braced"},
		{"concated literals", `value1 # value2,`, "0005_concated"},
		{"concated quote and literal", `"value1" # value2,`, "0006_concat_quote_literal"},
		// key = value
		{"simple name", `name = value,`, "0007_key_value"},
		{"simple name (compact)", `name=value,`, "0008_simple_name_compact"},
		{"name + compact value", `name=a#"b",`, "0009_name_compact_value"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// read the tag
			var tag BibTag
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibtag_read", tt.asset+".json"), &tag)
			// write the buffer
			writer := &bytes.Buffer{}
			if err := tag.Write(writer); (err != nil) != false {
				t.Errorf("BibTag.Write() error = %v, wantErr %v", err, false)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantString {
				t.Errorf("BibTag.Write() = %v, want %v", gotWriter, tt.wantString)
			}
		})
	}
}

func TestBibTag_Empty(t *testing.T) {
	tests := []struct {
		name  string
		asset string
		want  bool
	}{
		{"empty tag", "0001_empty", true},
		{"endingbrace", "0010_endingbrace", true},
		{"literal value", "0002_literal", false},
		{"quoted value", "0003_quoted", false},
		{"braced value", "0004_braced", false},
		{"concated literals", "0005_concated", false},
		{"concated quote and literal", "0006_concat_quote_literal", false},
		// key = value
		{"simple name", "0007_key_value", false},
		{"simple name (compact)", "0008_simple_name_compact", false},
		{"name + compact value", "0009_name_compact_value", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// read the tag
			var tag BibTag
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibtag_read", tt.asset+".json"), &tag)
			// call .Empty()
			if got := tag.Empty(); got != tt.want {
				t.Errorf("BibTag.Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBibTag_IsKeyValue(t *testing.T) {
	tests := []struct {
		name  string
		asset string
		want  bool
	}{
		{"empty tag", "0001_empty", false},
		{"endingbrace", "0010_endingbrace", false},
		{"literal value", "0002_literal", false},
		{"quoted value", "0003_quoted", false},
		{"braced value", "0004_braced", false},
		{"concated literals", "0005_concated", false},
		{"concated quote and literal", "0006_concat_quote_literal", false},
		// key = value
		{"simple name", "0007_key_value", true},
		{"simple name (compact)", "0008_simple_name_compact", true},
		{"name + compact value", "0009_name_compact_value", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// read the tag
			var tag BibTag
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibtag_read", tt.asset+".json"), &tag)
			// call .IsKeyValue()
			if got := tag.IsKeyValue(); got != tt.want {
				t.Errorf("BibTag.IsKeyValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBibTag_GetKey(t *testing.T) {
	tests := []struct {
		name  string
		asset string
	}{
		{"empty tag", "0001_empty"},
		{"endingbrace", "0010_endingbrace"},
		{"literal value", "0002_literal"},
		{"quoted value", "0003_quoted"},
		{"braced value", "0004_braced"},
		{"concated literals", "0005_concated"},
		{"concated quote and literal", "0006_concat_quote_literal"},
		// key = value
		{"simple name", "0007_key_value"},
		{"simple name (compact)", "0008_simple_name_compact"},
		{"name + compact value", "0009_name_compact_value"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// read the tag
			var tag BibTag
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibtag_read", tt.asset+".json"), &tag)
			// read the element
			var want *BibTagElement
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibtag_getkey", tt.asset+".json"), &want)
			// call GetKey
			if got := tag.GetKey(); !reflect.DeepEqual(got, want) {
				t.Errorf("BibTag.GetKey() = %v, want %v", got, want)
			}
		})
	}
}

func TestBibTag_GetValue(t *testing.T) {
	tests := []struct {
		name  string
		asset string
	}{
		{"empty tag", "0001_empty"},
		{"endingbrace", "0010_endingbrace"},
		{"literal value", "0002_literal"},
		{"quoted value", "0003_quoted"},
		{"braced value", "0004_braced"},
		{"concated literals", "0005_concated"},
		{"concated quote and literal", "0006_concat_quote_literal"},
		// key = value
		{"simple name", "0007_key_value"},
		{"simple name (compact)", "0008_simple_name_compact"},
		{"name + compact value", "0009_name_compact_value"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// read the tag
			var tag BibTag
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibtag_read", tt.asset+".json"), &tag)
			// read the elements
			var want []*BibTagElement
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibtag_getvalue", tt.asset+".json"), &want)
			if got := tag.GetValue(); !reflect.DeepEqual(got, want) {
				t.Errorf("BibTag.GetValue() = %v, want %v", got, want)
			}
		})
	}
}
