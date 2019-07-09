package bibliography

import (
	"reflect"
	"testing"

	"github.com/tkw1536/gotexml/utils"
)

func TestBibFile_readQuote(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantQuote BibString
	}{
		{
			"empty quotes",
			`""`,
			BibString{
				kind:  BibStringQuote,
				value: ``,
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 2, EOF: false},
				},
			},
		},
		{
			"simple quote",
			`"hello"`,
			BibString{
				kind:  BibStringQuote,
				value: `hello`,
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 7, EOF: false},
				},
			},
		},
		{
			"with { s",
			`"{\"}"`,
			BibString{
				kind:  BibStringQuote,
				value: `{\"}`,
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 6, EOF: false},
				},
			},
		},
		{
			"quote with spaces",
			`"hello world"`,
			BibString{
				kind:  BibStringQuote,
				value: `hello world`,
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 13, EOF: false},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := &BibFile{
				Reader: utils.NewRuneReaderFromString(tt.input + " , "),
			}
			gotQuote, err := bf.readQuote()
			if (err != nil) != false {
				t.Errorf("BibFile.readQuote() error = %v, wantErr %v", err, false)
				return
			}
			if !reflect.DeepEqual(gotQuote, tt.wantQuote) {
				t.Errorf("BibFile.readQuote() = %v, want %v", gotQuote, tt.wantQuote)
			}
		})
	}
}

func TestBibFile_readBrace(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantBrace BibString
	}{
		{
			"empty braces",
			`{}`,
			BibString{
				kind:  BibStringBracket,
				value: ``,
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 2, EOF: false},
				},
			},
		},
		{
			"simple braces",
			`{hello}`,
			BibString{
				kind:  BibStringBracket,
				value: `hello`,
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 7, EOF: false},
				},
			},
		},
		{
			"nested braces",
			`{hello{world}}`,
			BibString{
				kind:  BibStringBracket,
				value: `hello{world}`,
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 14, EOF: false},
				},
			},
		},
		{
			"brace with open \\",
			`{hello \{world}}`,
			BibString{
				kind:  BibStringBracket,
				value: `hello \{world}`,
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 16, EOF: false},
				},
			},
		},
		{
			"brace with close \\",
			`{hello world\}}`,
			BibString{
				kind:  BibStringBracket,
				value: `hello world\`,
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 14, EOF: false},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := &BibFile{
				Reader: utils.NewRuneReaderFromString(tt.input),
			}
			gotBrace, err := bf.readBrace()
			if (err != nil) != false {
				t.Errorf("BibFile.readBrace() error = %v, wantErr %v", err, false)
				return
			}
			if !reflect.DeepEqual(gotBrace, tt.wantBrace) {
				t.Errorf("BibFile.readBrace() = %v, want %v", gotBrace, tt.wantBrace)
			}
		})
	}
}

func TestBibFile_readLiteral(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantLit   BibString
		wantSpace BibString
	}{
		{
			"empty",
			`,`,
			BibString{
				kind:  BibStringLiteral,
				value: ``,
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
				},
			},
			BibString{
				kind:  BibStringOther,
				value: ``,
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
				},
			},
		},
		{
			"space",
			`hello world`,
			BibString{
				kind:  BibStringLiteral,
				value: `hello world`,
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 11, EOF: false},
				},
			},
			BibString{
				kind:  BibStringOther,
				value: ``,
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 11, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 11, EOF: false},
				},
			},
		},
		{
			"with an @ sign",
			`hello@world`,
			BibString{
				kind:  BibStringLiteral,
				value: `hello@world`,
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 11, EOF: false},
				},
			},
			BibString{
				kind:  BibStringOther,
				value: ``,
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 11, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 11, EOF: false},
				},
			},
		},
		{
			"with an \" sign",
			`hello"world`,
			BibString{
				kind:  BibStringLiteral,
				value: `hello"world`,
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 11, EOF: false},
				},
			},
			BibString{
				kind:  BibStringOther,
				value: ``,
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 11, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 11, EOF: false},
				},
			},
		},
		{
			"surrounding space",
			`hello  world     `,
			BibString{
				kind:  BibStringLiteral,
				value: `hello  world`,
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 12, EOF: false},
				},
			},
			BibString{
				kind:  BibStringOther,
				value: `     `,
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 12, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 16, EOF: false},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := &BibFile{
				Reader: utils.NewRuneReaderFromString(tt.input + "},"),
			}
			gotLit, gotSpace, err := bf.readLiteral()
			if (err != nil) != false {
				t.Errorf("BibFile.readLiteral() error = %v, wantErr %v", err, false)
				return
			}
			if !reflect.DeepEqual(gotLit, tt.wantLit) {
				t.Errorf("BibFile.readLiteral() gotLit = %v, wantLit %v", gotLit, tt.wantLit)
			}
			if !reflect.DeepEqual(gotSpace, tt.wantSpace) {
				t.Errorf("BibFile.readLiteral() gotSpace = %v, wantSpace %v", gotSpace, tt.wantSpace)
			}
		})
	}
}
