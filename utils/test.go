package utils

import (
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"os"
)

// ReadFileOrPanic reads a string from a file or panics
func ReadFileOrPanic(filename string) string {
	if filename == "" {
		panic("missing filename")
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

// UnmarshalFileOrPanic json.Unmarshals data from a file or panic()s
func UnmarshalFileOrPanic(filename string, data interface{}) {
	if filename == "" {
		panic("Missing filename")
	}
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(bytes, data); err != nil {
		panic(err)
	}
}

// CompressUnmarshalFileOrPanic uncompresses and json.Unmarshals data from a file or panic()s
func CompressUnmarshalFileOrPanic(filename string, data interface{}) {
	if filename == "" {
		panic("Missing filename")
	}

	// opent the file for reading
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// prepare unpacking
	z, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	defer z.Close()

	// and read it all
	bytes, err := ioutil.ReadAll(z)
	if err != nil {
		panic(err)
	}

	// finally unmarshal
	if err = json.Unmarshal(bytes, data); err != nil {
		panic(err)
	}
}

// MarshalFileOrPanic marshals data into a file or panics
func MarshalFileOrPanic(filename string, data interface{}) {
	if filename == "" {
		panic("Missing filename")
	}
	bytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(filename, bytes, 0644); err != nil {
		panic(err)
	}
}

// CompressMarshalFileOrPanic marshals and compresses data into a file or panics
func CompressMarshalFileOrPanic(filename string, data interface{}) {
	if filename == "" {
		panic("Missing filename")
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// marshal the data
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	// and write it out to the file
	w := gzip.NewWriter(f)
	w.Write(bytes)
	w.Close()
}
