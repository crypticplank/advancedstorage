package advancedstorage

import (
	"bytes"
	"encoding/gob"
	"log"
	"testing"
)

type TestStruct struct {
	TestText   string
	TestData   []byte
	TestNumber int64
}

func TestStorage_WriteToFile(t *testing.T) {
	var buffer bytes.Buffer
	s, _ := New(&Options{Filename: "secrets", Verbose: true})
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(TestStruct{TestText: "testing", TestNumber: 420, TestData: []byte{0x41, 0x41}})
	s.WriteToFile(buffer.Bytes())
}

func TestStorage_ReadFromFile(t *testing.T) {
	s, _ := New(&Options{Filename: "secrets", Verbose: true})
	data, _ := s.ReadFromFile()
	var test TestStruct
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	decoder.Decode(&test)
	log.Println(test)
}
