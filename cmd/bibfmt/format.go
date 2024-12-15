package main

import (
	"io"
	"os"

	"github.com/tkw1536/gotexml/bibliography"
	"github.com/tkw1536/gotexml/utils"
)

func main() {
	argc := len(os.Args) - 1
	argv := os.Args[1:]

	// read from the argv[0] or stdin
	var inReader io.Reader
	if argc == 0 {
		inReader = os.Stdin
	} else {
		var err error
		var f *os.File
		if f, err = os.Open(argv[0]); err != nil {
			panic(err)
		}
		inReader = f
		defer f.Close()
	}

	// write to argv[1] or stdout
	var outWriter io.Writer
	if argc < 2 {
		outWriter = os.Stdout
	} else {
		var err error
		var f *os.File
		if f, err = os.OpenFile(argv[1], os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			panic(err)
		}
		outWriter = f
		defer f.Close()
	}

	// and start reading
	reader := utils.NewRuneReaderFromReader(inReader)
	file, err := bibliography.NewBibFileFromReader(reader)
	if err != nil {
		panic(err)
	}

	// format the file using the default formatter
	bibliography.DefaultFormatter.Format(file)

	// and write out the formatted file
	if err := file.Write(outWriter); err != nil {
		panic(err)
	}
}
