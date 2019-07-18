package bibliography

import (
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/tkw1536/gotexml/utils"
)

func Test_readFile(t *testing.T) {
	tests := []struct {
		name  string
		asset string
	}{
		{"complicated.bib", "0001_complicated"},
		{"kwarc.bib", "0002_kwarc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// read input
			file, err := os.Open(path.Join("testdata", "bibfile_read", tt.asset+".bib"))
			if err != nil {
				panic(err)
			}
			defer file.Close()

			// call readEntry
			gotFile, err := readFile(utils.NewRuneReaderFromReader(file))

			// read the assets
			var wantFile BibFile
			utils.CompressUnmarshalFileOrPanic(path.Join("testdata", "bibfile_read", tt.asset+".json.gz"), &wantFile)

			if (err != nil) != false {
				t.Errorf("BibTag.readFile() error = %v, wantErr %v", err, false)
				return
			}

			if !reflect.DeepEqual(gotFile, wantFile) {
				t.Errorf("BibTag.readFile() = %v, want %v", gotFile, wantFile)
			}
		})
	}
}

func Benchmark_ReadFile_Complicated(b *testing.B) {
	benchmarkReadFile(complicatedBibFileText, b)
}

func Benchmark_ReadFile_Kwarc(b *testing.B) {
	benchmarkReadFile(kwarcBibFileText, b)
}

func benchmarkReadFile(content string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		readFile(utils.NewRuneReaderFromString(content))
	}
}

var complicatedBibFileText = utils.ReadFileOrPanic(path.Join("testdata", "bibfile_read", "0001_complicated.bib"))
var kwarcBibFileText = utils.ReadFileOrPanic(path.Join("testdata", "bibfile_read", "0002_kwarc.bib"))
