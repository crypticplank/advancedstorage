package advancedstorage

import (
	"log"
	"testing"
)

func TestStorage_WriteToFile(t *testing.T) {
	s, _ := New(&Options{Filename: "secrets"})
	s.WriteToFile([]byte("sheeeeeeeeesh"))
}

func TestStorage_ReadFromFile(t *testing.T) {
	s, _ := New(&Options{Filename: "secrets"})
	data, _ := s.ReadFromFile()
	log.Println(string(data))
}
