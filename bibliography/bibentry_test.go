package bibliography

import (
	"bytes"
	"io"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/tkw1536/gotexml/utils"
)

func Test_readEntry(t *testing.T) {
	tests := []struct {
		name    string
		asset   string
		wantEOF bool
	}{
		{"empty", "0001_empty", true},
		{"preamble entry", "0002_preamble", false},
		{"string entry", "0003_string", false},
		{"inproceedings entry", "0004_inproceedings", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// read input
			file, err := os.Open(path.Join("testdata", "bibentry_read", tt.asset+".bib"))
			if err != nil {
				panic(err)
			}
			defer file.Close()

			// call readEntry
			gotEntry := &BibEntry{}
			err = gotEntry.readEntry(utils.NewRuneReaderFromReader(file))

			// read the assets
			var wantEntry *BibEntry
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibentry_read", tt.asset+".json"), &wantEntry)

			// if we want eof, only test for EOF
			if tt.wantEOF {
				if err != io.EOF {
					t.Errorf("BibField.readEntry() error = %v, wantErr %v", err, io.EOF)
					return
				}
				return
			}

			if (err != nil) != false {
				t.Errorf("BibField.readEntry() error = %v, wantErr %v", err, false)
				return
			}

			if !reflect.DeepEqual(gotEntry, wantEntry) {
				t.Errorf("BibField.readEntry() = %v, want %v", gotEntry, wantEntry)
			}
		})
	}
}

func Benchmark_ReadEntry_Empty(b *testing.B) {
	benchmarkReadEntry(emptyEntryText, b)
}

func Benchmark_ReadEntry_Preamble(b *testing.B) {
	benchmarkReadEntry(preambleEntryText, b)
}
func Benchmark_ReadEntry_String(b *testing.B) {
	benchmarkReadEntry(stringEntryText, b)
}
func Benchmark_ReadEntry_Inproceedings(b *testing.B) {
	benchmarkReadEntry(inproceedingsEntryText, b)
}

func benchmarkReadEntry(content string, b *testing.B) {
	entry := &BibEntry{}
	for n := 0; n < b.N; n++ {
		entry.readEntry(utils.NewRuneReaderFromString(content))
	}
}

var emptyEntryText = utils.ReadFileOrPanic(path.Join("testdata", "bibentry_read", "0001_empty.bib"))
var preambleEntryText = utils.ReadFileOrPanic(path.Join("testdata", "bibentry_read", "0002_preamble.bib"))
var stringEntryText = utils.ReadFileOrPanic(path.Join("testdata", "bibentry_read", "0003_string.bib"))
var inproceedingsEntryText = utils.ReadFileOrPanic(path.Join("testdata", "bibentry_read", "0004_inproceedings.bib"))

func TestBibEntry_Write(t *testing.T) {
	tests := []struct {
		name  string
		asset string
	}{
		{"preamble entry", "0002_preamble"},
		{"string entry", "0003_string"},
		{"inproceedings entry", "0004_inproceedings"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// load the entry
			var entry BibEntry
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibentry_read", tt.asset+".json"), &entry)
			// load the string we want
			wantString := utils.ReadFileOrPanic(path.Join("testdata", "bibentry_read", tt.asset+".bib"))
			// write the buffer
			writer := &bytes.Buffer{}
			if err := entry.Write(writer); (err != nil) != false {
				t.Errorf("BibEntry.Write() error = %v, wantErr %v", err, false)
				return
			}
			if gotWriter := writer.String(); gotWriter != wantString {
				t.Errorf("BibEntry.Write() = %q, want %q", gotWriter, wantString)
			}
		})
	}
}
