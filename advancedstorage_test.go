package advancedstorage

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"testing"
)

type TestStruct struct {
	TestText   string
	TestData   []byte
	TestNumber int64
}

var (
	Test = TestStruct{TestText: "testing", TestNumber: 420, TestData: []byte{0x41, 0x41}}
)

func TestStorage_WriteToFile(t *testing.T) {
	var buffer bytes.Buffer
	s, _ := New(&Options{Filename: "advancedstorage.test"})
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(Test)
	s.WriteToFile(buffer.Bytes())
	// io.ReadAll(s.Reader)
}

func TestStorage_ReadFromFile(t *testing.T) {
	s, _ := New(&Options{Filename: "advancedstorage.test"})
	data, _ := s.ReadFromFile()
	var test TestStruct
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	decoder.Decode(&test)
	if !reflect.DeepEqual(test, Test) {
		t.Fatalf("Data could not be verified %s", t.Name())
	}
}
