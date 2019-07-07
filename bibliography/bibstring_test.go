package bibliography

import (
	"reflect"
	"testing"

	"github.com/tkw1536/gotexml/utils"
)

func TestBibString_Copy(t *testing.T) {
	type fields struct {
		kind   BibStringKind
		value  string
		source utils.ReaderRange
	}
	tests := []struct {
		name   string
		fields fields
		want   *BibString
	}{
		{"copy a bibstring returns a copy", fields{
			kind:  BibStringOther,
			value: "something",
			source: utils.ReaderRange{
				Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
				End:   utils.ReaderPosition{Line: 2, Column: 1, EOF: true},
			},
		}, &BibString{
			kind:  BibStringOther,
			value: "something",
			source: utils.ReaderRange{
				Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
				End:   utils.ReaderPosition{Line: 2, Column: 1, EOF: true},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := &BibString{
				kind:   tt.fields.kind,
				value:  tt.fields.value,
				source: tt.fields.source,
			}
			if got := bs.Copy(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BibString.Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBibString_NormalizeValue(t *testing.T) {
	tests := []struct {
		name   string
		before *BibString
		after  *BibString
	}{
		{"normalize simple value", &BibString{value: "HeLlO wOrLd"}, &BibString{value: "hello world"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before.NormalizeValue()
			if !reflect.DeepEqual(tt.before, tt.after) {
				t.Errorf("BibString.NormalizeValue() = %v, want %v", tt.before, tt.after)
			}
		})
	}
}

func TestBibString_Evaluate(t *testing.T) {
	tests := []struct {
		name    string
		before  *BibString
		context map[string]string

		wantOK bool
		after  *BibString
	}{
		{"evaluating other bibstring", &BibString{
			kind:  BibStringOther,
			value: "something",
			source: utils.ReaderRange{
				Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
				End:   utils.ReaderPosition{Line: 2, Column: 0, EOF: true},
			},
		}, map[string]string{"something": "other"}, true, &BibString{
			kind:  BibStringOther,
			value: "something",
			source: utils.ReaderRange{
				Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
				End:   utils.ReaderPosition{Line: 2, Column: 0, EOF: true},
			},
		}},

		{"evaluating bracket bibstring", &BibString{
			kind:  BibStringBracket,
			value: "something",
			source: utils.ReaderRange{
				Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
				End:   utils.ReaderPosition{Line: 2, Column: 0, EOF: true},
			},
		}, map[string]string{"something": "other"}, true, &BibString{
			kind:  BibStringBracket,
			value: "something",
			source: utils.ReaderRange{
				Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
				End:   utils.ReaderPosition{Line: 2, Column: 0, EOF: true},
			},
		}},

		{"evaluating literal bibstring with valid context", &BibString{
			kind:  BibStringLiteral,
			value: "something",
			source: utils.ReaderRange{
				Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
				End:   utils.ReaderPosition{Line: 2, Column: 0, EOF: true},
			}}, map[string]string{"something": "other"}, true, &BibString{
			kind:  BibStringEvaluated,
			value: "other",
			source: utils.ReaderRange{
				Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
				End:   utils.ReaderPosition{Line: 2, Column: 0, EOF: true},
			},
		}},

		{"evaluating literal bibstring with invalid context", &BibString{
			kind:  BibStringLiteral,
			value: "something",
			source: utils.ReaderRange{
				Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
				End:   utils.ReaderPosition{Line: 2, Column: 0, EOF: true},
			},
		}, nil, false, &BibString{
			kind:  BibStringLiteral,
			value: "something",
			source: utils.ReaderRange{
				Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
				End:   utils.ReaderPosition{Line: 2, Column: 0, EOF: true},
			},
		}},

		{"evaluating quoted bibstring", &BibString{
			kind:  BibStringQuote,
			value: "something",
			source: utils.ReaderRange{
				Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
				End:   utils.ReaderPosition{Line: 2, Column: 0, EOF: true},
			},
		}, map[string]string{"something": "other"}, true, &BibString{
			kind:  BibStringQuote,
			value: "something",
			source: utils.ReaderRange{
				Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
				End:   utils.ReaderPosition{Line: 2, Column: 0, EOF: true},
			},
		}},

		{"evaluating evaluated bibstring", &BibString{
			kind:  BibStringEvaluated,
			value: "something",
			source: utils.ReaderRange{
				Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
				End:   utils.ReaderPosition{Line: 2, Column: 0, EOF: true},
			},
		}, map[string]string{"something": "other"}, true, &BibString{
			kind:  BibStringEvaluated,
			value: "something",
			source: utils.ReaderRange{
				Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
				End:   utils.ReaderPosition{Line: 2, Column: 0, EOF: true},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.before.Evaluate(tt.context); got != tt.wantOK {
				t.Errorf("BibString.Evaluate() ok = %v, wantOK %v", got, tt.wantOK)
			}
			if !reflect.DeepEqual(tt.before, tt.after) {
				t.Errorf("BibString.Evaluate() = %v, want %v", tt.before, tt.after)
			}
		})
	}
}

func TestBibString_Append(t *testing.T) {
	tests := []struct {
		name   string
		before *BibString
		other  *BibString

		after *BibString
	}{
		{"adding two strings", &BibString{
			kind:  BibStringQuote,
			value: "hello \n",
			source: utils.ReaderRange{
				Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
				End:   utils.ReaderPosition{Line: 2, Column: 0, EOF: false},
			}}, &BibString{
			kind:  BibStringQuote,
			value: "world\n",
			source: utils.ReaderRange{
				Start: utils.ReaderPosition{Line: 2, Column: 1, EOF: false},
				End:   utils.ReaderPosition{Line: 3, Column: 0, EOF: true},
			}}, &BibString{
			kind:  BibStringEvaluated,
			value: "hello \nworld\n",
			source: utils.ReaderRange{
				Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
				End:   utils.ReaderPosition{Line: 3, Column: 0, EOF: true},
			}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before.Append(tt.other)
			if !reflect.DeepEqual(tt.before, tt.after) {
				t.Errorf("BibString.Append() = %v, want %v", tt.before, tt.after)
			}
		})
	}
}
