package bibliography

import (
	"bytes"
	"path"
	"reflect"
	"testing"

	"github.com/tkw1536/gotexml/utils"
)

func Test_readField(t *testing.T) {
	tests := []struct {
		name  string
		input string
		asset string
	}{
		// value only
		{"empty field", ``, "0001_empty"},
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
			var wantField *BibField
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibfield_read", tt.asset+".json"), &wantField)

			// call readField
			gotField := &BibField{}
			err := gotField.readField(utils.NewRuneReaderFromString(tt.input + ", "))

			if (err != nil) != false {
				t.Errorf("BibField.readField() error = %v, wantErr %v", err, false)
				return
			}

			if !reflect.DeepEqual(gotField, wantField) {
				t.Errorf("BibField.readField() = %v, want %v", gotField, wantField)
			}
		})
	}
}

func Benchmark_ReadField_Empty(b *testing.B) {
	benchmarkReadField(``, b)
}

func Benchmark_ReadField_Value(b *testing.B) {
	benchmarkReadField(`value`, b)
}

func Benchmark_ReadField_QuotedValue(b *testing.B) {
	benchmarkReadField(`"value"`, b)
}

func Benchmark_ReadField_BracketedValue(b *testing.B) {
	benchmarkReadField(`{value}`, b)
}

func Benchmark_ReadField_Concat(b *testing.B) {
	benchmarkReadField(`value1 # value2`, b)
}

func Benchmark_ReadField_ComplexConcat(b *testing.B) {
	benchmarkReadField(`"value1" # value2`, b)
}

func Benchmark_ReadField_KeyValue(b *testing.B) {
	benchmarkReadField(`name = value`, b)
}

func Benchmark_ReadField_KeyValueCompact(b *testing.B) {
	benchmarkReadField(`name=value`, b)
}

func Benchmark_ReadField_KeyValueComplex(b *testing.B) {
	benchmarkReadField(`name=a#"b"`, b)
}

func benchmarkReadField(content string, b *testing.B) {
	p := content + ", "
	field := &BibField{}
	for n := 0; n < b.N; n++ {
		field.readField(utils.NewRuneReaderFromString(p))
	}
}

func TestBibField_Write(t *testing.T) {
	tests := []struct {
		name       string
		wantString string
		asset      string
	}{
		// value only
		{"empty field", `,`, "0001_empty"},
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
			// read the field
			var field BibField
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibfield_read", tt.asset+".json"), &field)
			// write the buffer
			writer := &bytes.Buffer{}
			if err := field.Write(writer); (err != nil) != false {
				t.Errorf("BibField.Write() error = %v, wantErr %v", err, false)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantString {
				t.Errorf("BibField.Write() = %v, want %v", gotWriter, tt.wantString)
			}
		})
	}
}

func TestBibField_Empty(t *testing.T) {
	tests := []struct {
		name  string
		asset string
		want  bool
	}{
		{"empty field", "0001_empty", true},
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
			// read the field
			var field BibField
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibfield_read", tt.asset+".json"), &field)
			// call .Empty()
			if got := field.Empty(); got != tt.want {
				t.Errorf("BibField.Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBibField_IsKeyValue(t *testing.T) {
	tests := []struct {
		name  string
		asset string
		want  bool
	}{
		{"empty field", "0001_empty", false},
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
			// read the field
			var field BibField
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibfield_read", tt.asset+".json"), &field)
			// call .IsKeyValue()
			if got := field.IsKeyValue(); got != tt.want {
				t.Errorf("BibField.IsKeyValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBibField_GetKey(t *testing.T) {
	tests := []struct {
		name  string
		asset string
	}{
		{"empty field", "0001_empty"},
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
			// read the field
			var field BibField
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibfield_read", tt.asset+".json"), &field)
			// read the element
			var want *BibFieldElement
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibfield_getkey", tt.asset+".json"), &want)
			// call GetKey
			if got := field.GetKey(); !reflect.DeepEqual(got, want) {
				t.Errorf("BibField.GetKey() = %v, want %v", got, want)
			}
		})
	}
}

func TestBibField_GetValue(t *testing.T) {
	tests := []struct {
		name  string
		asset string
	}{
		{"empty field", "0001_empty"},
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
			// read the field
			var field BibField
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibfield_read", tt.asset+".json"), &field)
			// read the elements
			var want []*BibFieldElement
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibfield_getvalue", tt.asset+".json"), &want)
			if got := field.GetValue(); !reflect.DeepEqual(got, want) {
				t.Errorf("BibField.GetValue() = %v, want %v", got, want)
			}
		})
	}
}
