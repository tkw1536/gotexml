package bibliography

import (
	"io"
	"path"
	"reflect"
	"testing"

	"github.com/tkw1536/gotexml/utils"
)

func Test_readTag(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		asset   string
		wantEOF bool
	}{
		// value only
		{"empty tag", ``, "0001_empty", true},
		{"literal value", `value`, "0002_literal", false},
		{"quoted value", `"value"`, "0003_quoted", false},
		{"braced value", `{value}`, "0004_braced", false},
		{"concated literals", `value1 # value2`, "0005_concated", false},
		{"concated quote and literal", `"value1" # value2`, "0006_concat_quote_literal", false},
		// key = value
		{"simple name", `name = value`, "0007_key_value", false},
		{"simple name (compact)", `name=value`, "0008_simple_name_compact", false},
		{"name + compact value", `name=a#"b"`, "0009_name_compact_value", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// read the assets
			var wantTag BibTag
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibtag_read", tt.asset+".json"), &wantTag)

			// call readTag
			gotTag, err := readTag(utils.NewRuneReaderFromString(tt.input + ", "))

			// if we want eof, only test for EOF
			if tt.wantEOF {
				if err != io.EOF {
					t.Errorf("BibTag.readTag() error = %v, wantErr %v", err, io.EOF)
					return
				}
				return
			}

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
