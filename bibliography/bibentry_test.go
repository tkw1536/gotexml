package bibliography

import (
	"io"
	"io/ioutil"
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

			// read the assets
			var wantEntry BibEntry
			utils.UnmarshalFileOrPanic(path.Join("testdata", "bibentry_read", tt.asset+".json"), &wantEntry)

			// call readEntry
			gotEntry, err := readEntry(utils.NewRuneReaderFromReader(file))

			// if we want eof, only test for EOF
			if tt.wantEOF {
				if err != io.EOF {
					t.Errorf("BibTag.readEntry() error = %v, wantErr %v", err, io.EOF)
					return
				}
				return
			}

			if (err != nil) != false {
				t.Errorf("BibTag.readEntry() error = %v, wantErr %v", err, false)
				return
			}

			if !reflect.DeepEqual(gotEntry, wantEntry) {
				t.Errorf("BibTag.readEntry() = %v, want %v", gotEntry, wantEntry)
			}
		})
	}
}

func Benchmark_ReadEntry(b *testing.B) {

	for n := 0; n < b.N; n++ {
		for _, a := range benchmarkReadEntryAssets {
			readEntry(utils.NewRuneReaderFromString(a))
		}
	}
}

var benchmarkReadEntryAssets []string

func init() {
	// read in all the assets for the benchmark
	filenames := []string{"0001_empty", "0002_preamble", "0003_string", "0004_inproceedings"}
	benchmarkReadEntryAssets = make([]string, len(filenames))

	var err error
	var d []byte
	for idx, asset := range filenames {
		d, err = ioutil.ReadFile(path.Join("testdata", "bibentry_read", asset+".json"))
		if err != nil {
			panic(err)
		}
		benchmarkReadEntryAssets[idx] = string(d)
	}

}
