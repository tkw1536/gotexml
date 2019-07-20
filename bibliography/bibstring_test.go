package bibliography

import (
	"bytes"
	"path"
	"reflect"
	"testing"

	"github.com/tkw1536/gotexml/utils"
)

func TestBibString_Append(t *testing.T) {
	tests := []struct {
		name  string
		asset string
	}{
		{"adding two strings", "0001_adding_two_strings"},
		{"adding empty string", "0002_adding_empty_string"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// read the assets
			var before, other, after BibString
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibstring_append", tt.asset+"_before.json"), &before)
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibstring_append", tt.asset+"_other.json"), &other)
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibstring_append", tt.asset+"_after.json"), &after)

			// call append
			before.Append(&other)
			if !reflect.DeepEqual(before, after) {
				t.Errorf("BibString.Append() = %v, want %v", before, after)
			}
		})
	}
}

func TestBibString_AppendRaw(t *testing.T) {
	tests := []struct {
		name  string
		value string
		asset string
	}{
		{"adding two strings", "world\n", "0001_adding_two_strings"},
		{"adding empty string", "", "0002_adding_empty_string"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// read the others
			var loc utils.ReaderRange
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibstring_appendraw", tt.asset+"_loc.json"), &loc)

			var before, after BibString
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibstring_appendraw", tt.asset+"_before.json"), &before)
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibstring_appendraw", tt.asset+"_after.json"), &after)

			// call appendraw
			before.AppendRaw(tt.value, loc)
			if !reflect.DeepEqual(before, after) {
				t.Errorf("BibString.AppendRaw() = %v, want %v", before, after)
			}
		})
	}
}

func TestBibString_readQuote(t *testing.T) {
	tests := []struct {
		name  string
		input string
		asset string
	}{
		{"empty quotes", `""`, "0001_empty"},
		{"simple quote", `"hello"`, "0002_simple_quote"},
		{"with { s", `"{\"}"`, "0003_with_curly"},
		{"quote with spaces", `"hello world"`, "0004_with_spaces"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// read the asset
			var wantQuote BibString
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibstring_quote", tt.asset+".json"), &wantQuote)

			// call readQuote
			gotQuote := &BibString{}
			err := gotQuote.readQuote(utils.NewRuneReaderFromString(tt.input + " , "))
			if (err != nil) != false {
				t.Errorf("BibString.readQuote() error = %v, wantErr %v", err, false)
				return
			}
			if !reflect.DeepEqual(gotQuote, &wantQuote) {
				t.Errorf("BibString.readQuote() = %v, want %v", gotQuote, &wantQuote)
			}
		})
	}
}

func Benchmark_ReadQuote_Empty(b *testing.B) {
	benchmarkReadQuote(`""`, b)
}
func Benchmark_ReadQuote_Simple(b *testing.B) {
	benchmarkReadQuote(`"hello"`, b)
}
func Benchmark_ReadQuote_WithCurly(b *testing.B) {
	benchmarkReadQuote(`"{\"}"`, b)
}
func Benchmark_ReadQuote_WithSpaces(b *testing.B) {
	benchmarkReadQuote(`"hello world"`, b)
}

func benchmarkReadQuote(content string, b *testing.B) {
	p := content + " , "
	bs := &BibString{}
	for n := 0; n < b.N; n++ {
		bs.readQuote(utils.NewRuneReaderFromString(p))
	}
}

func TestBibFile_readBrace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		asset string
	}{
		{"empty braces", `{}`, "0001_empty"},
		{"simple braces", `{hello}`, "0002_simple"},
		{"nested braces", `{hello{world}}`, "0003_nested"},
		{"brace with open \\", `{hello \{world}}`, "0004_open_slashes"},
		{"brace with close \\", `{hello world\}}`, "0005_close_slashes"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// read the asset
			var wantBrace BibString
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibstring_brace", tt.asset+".json"), &wantBrace)

			// call readBrace
			gotBrace := &BibString{}
			err := gotBrace.readBrace(utils.NewRuneReaderFromString(tt.input))
			if (err != nil) != false {
				t.Errorf("BibString.readBrace() error = %v, wantErr %v", err, false)
				return
			}
			if !reflect.DeepEqual(gotBrace, &wantBrace) {
				t.Errorf("BibString.readBrace() = %v, want %v", gotBrace, &wantBrace)
			}
		})
	}
}

func Benchmark_ReadBrace_Empty(b *testing.B) {
	benchmarkReadBrace(`{}`, b)
}

func Benchmark_ReadBrace_Simple(b *testing.B) {
	benchmarkReadBrace(`{hello}`, b)
}

func Benchmark_ReadBrace_Nested(b *testing.B) {
	benchmarkReadBrace(`{hello{world}}`, b)
}

func Benchmark_ReadBrace_OpenSlashes(b *testing.B) {
	benchmarkReadBrace(`{hello \{world}}`, b)
}

func Benchmark_ReadBrace_CloseSlashes(b *testing.B) {
	benchmarkReadBrace(`{hello world\}}`, b)
}

func benchmarkReadBrace(content string, b *testing.B) {
	bs := &BibString{}
	for n := 0; n < b.N; n++ {
		bs.readBrace(utils.NewRuneReaderFromString(content))
	}
}

func TestBibFile_readLiteral(t *testing.T) {
	tests := []struct {
		name  string
		input string
		asset string
	}{
		{"empty", `,`, "0001_empty"},
		{"one character", `a`, "0002_one_character"},
		{"space", `hello world`, "0003_space"},
		{"with an @ sign", `hello@world`, "0004_with_at_sign"},
		{"with an \" sign", `hello"world`, "0005_with_quote_sign"},
		{"surrounding space", `hello  world     `, "0006_surrounding_space"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// read test assets
			var wantLit, wantSpace BibString
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibstring_literal", tt.asset+"_lit.json"), &wantLit)
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibstring_literal", tt.asset+"_space.json"), &wantSpace)

			// call readLiteral
			gotLit := &BibString{}
			gotSpace, err := gotLit.readLiteral(utils.NewRuneReaderFromString(tt.input + "},"))
			if (err != nil) != false {
				t.Errorf("BibString.readLiteral() error = %v, wantErr %v", err, false)
				return
			}
			if !reflect.DeepEqual(gotLit, &wantLit) {
				t.Errorf("BibString.readLiteral() gotLit = %v, wantLit %v", gotLit, &wantLit)
			}
			if !reflect.DeepEqual(gotSpace, &wantSpace) {
				t.Errorf("BibString.readLiteral() gotSpace = %v, wantSpace %v", gotSpace, &wantSpace)
			}
		})
	}
}

func Benchmark_ReadLiteral_Empty(b *testing.B) {
	benchmarkReadLiteral(`,`, b)
}
func Benchmark_ReadLiteral_OneCharacter(b *testing.B) {
	benchmarkReadLiteral(`a`, b)
}
func Benchmark_ReadLiteral_Space(b *testing.B) {
	benchmarkReadLiteral(`hello world`, b)
}
func Benchmark_ReadLiteral_WithAtSign(b *testing.B) {
	benchmarkReadLiteral(`hello@world`, b)
}
func Benchmark_ReadLiteral_withQuoteSign(b *testing.B) {
	benchmarkReadLiteral(`hello"world`, b)
}
func Benchmark_ReadLiteral_SurroundingSpace(b *testing.B) {
	benchmarkReadLiteral(`hello  world     `, b)
}

func benchmarkReadLiteral(content string, b *testing.B) {
	p := content + "},"
	bs := &BibString{}
	for n := 0; n < b.N; n++ {
		bs.readLiteral(utils.NewRuneReaderFromString(p))
	}
}

func TestBibString_Write(t *testing.T) {
	tests := []struct {
		name       string
		wantString string
		asset      string
	}{
		// Quote
		{"empty quotes", `""`, path.Join("bibstring_quote", "0001_empty")},
		{"simple quote", `"hello"`, path.Join("bibstring_quote", "0002_simple_quote")},
		{"with { s", `"{\"}"`, path.Join("bibstring_quote", "0003_with_curly")},
		{"quote with spaces", `"hello world"`, path.Join("bibstring_quote", "0004_with_spaces")},
		// Brace
		{"empty braces", `{}`, path.Join("bibstring_brace", "0001_empty")},
		{"simple braces", `{hello}`, path.Join("bibstring_brace", "0002_simple")},
		{"nested braces", `{hello{world}}`, path.Join("bibstring_brace", "0003_nested")},
		{"brace with open \\", `{hello \{world}}`, path.Join("bibstring_brace", "0004_open_slashes")},
		{"brace with close \\", `{hello world\}`, path.Join("bibstring_brace", "0005_close_slashes")},
		// Literal
		{"empty", ``, path.Join("bibstring_literal", "0001_empty_lit")},
		{"one character", `a`, path.Join("bibstring_literal", "0002_one_character_lit")},
		{"space", `hello world`, path.Join("bibstring_literal", "0003_space_lit")},
		{"with an @ sign", `hello@world`, path.Join("bibstring_literal", "0004_with_at_sign_lit")},
		{"with an \" sign", `hello"world`, path.Join("bibstring_literal", "0005_with_quote_sign_lit")},
		{"surrounding space", `hello  world`, path.Join("bibstring_literal", "0006_surrounding_space_lit")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// read the literal
			bs := &BibString{}
			utils.UnmarshalFileOrPanic(path.Join("testdata", tt.asset+".json"), &bs)

			// check if write into a buffer suceeds
			writer := &bytes.Buffer{}
			if err := bs.Write(writer); (err != nil) != false {
				t.Errorf("BibString.Write() error = %v, wantErr %v", err, false)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantString {
				t.Errorf("BibString.Write() = %v, want %v", gotWriter, tt.wantString)
			}
		})
	}
}
