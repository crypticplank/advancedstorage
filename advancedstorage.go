package advancedstorage

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
)

type Storage struct {
	Filename string
	Verbose  bool
	Reader   io.Reader
}

type Options struct {
	Filename string
	Verbose  bool
}

type Header struct {
	DataSize      int
	EncryptionKey []byte
}

func (s *Storage) DoesFileExist() bool {
	_, err := os.Stat(s.Filename)
	return !errors.Is(err, fs.ErrNotExist)
}

func (s *Storage) WriteToFile(b []byte) error {
	file, err := os.Create(s.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return err
	}

	if s.Verbose {
		log.Printf("Key: 0x%s", hex.EncodeToString(key))
	}

	header := Header{
		EncryptionKey: key,
		DataSize:      len(b),
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	encoder.Encode(header)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	if s.Verbose {
		log.Printf("Nonce: 0x%s", hex.EncodeToString(nonce))
	}

	b = aesGCM.Seal(nonce, nonce, b, nil)

	HeaderSize := make([]byte, 8)

	binary.LittleEndian.PutUint64(HeaderSize, uint64(buf.Len()))

	data := append(HeaderSize, append(buf.Bytes(), b...)...)

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) ReadFromFile() ([]byte, error) {
	data, err := os.ReadFile(s.Filename)
	if err != nil {
		return nil, err
	}

	gobSizeBytes := data[:8]
	gobSize := int64(binary.LittleEndian.Uint64(gobSizeBytes))

	var header Header

	data = data[8:]

	buf := bytes.NewBuffer(data[:gobSize])
	decoder := gob.NewDecoder(buf)
	decoder.Decode(&header)

	data = data[gobSize:]

	block, err := aes.NewCipher(header.EncryptionKey)
	if s.Verbose {
		log.Printf("Key: 0x%s", hex.EncodeToString(header.EncryptionKey))
	}

	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	nonce, cipher := data[:nonceSize], data[nonceSize:]
	if s.Verbose {
		log.Printf("Nonce: 0x%s", hex.EncodeToString(nonce))
	}

	data, _ = aesGCM.Open(nil, nonce, cipher, nil)

	return data, nil
}

func New(options *Options) (*Storage, error) {
	if len(options.Filename) < 1 {
		return nil, errors.New("filename must be greater than one")
	}
	return &Storage{
		Filename: options.Filename,
		Verbose:  options.Verbose,
	}, nil
}
