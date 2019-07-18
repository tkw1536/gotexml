package bibliography

import (
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
		{"endingbrace", `}`, "0010_name_compact_value"},
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
			var wantTag BibTag
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibtag_read", tt.asset+".json"), &wantTag)

			// call readTag
			gotTag, err := readTag(utils.NewRuneReaderFromString(tt.input + ", "))

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
	for n := 0; n < b.N; n++ {
		readTag(utils.NewRuneReaderFromString(p))
	}
}
