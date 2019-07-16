package utils

import (
	"encoding/json"
	"io/ioutil"
)

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
