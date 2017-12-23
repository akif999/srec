package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"log"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/akif999/srec"
)

var (
	filename     = kingpin.Arg("Filename", "Srec filename").ExistingFile()
	startAddr    = kingpin.Arg("StartAddr", "address of start encryption").Uint32()
	sizeOfBlocks = kingpin.Arg("sizeOfBlocks", "Size of blocks of encryption").Uint32()

	key = []byte{
		0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
		0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
	}
	iv = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
)

func main() {
	sr := srec.NewSrec()

	kingpin.Parse()
	fmt.Printf("%08X\n", *startAddr)
	fmt.Printf("%08X\n", *sizeOfBlocks)

	fp, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	err = sr.Parse(fp)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(sr.String())
	bytes := getBytes(sr, *startAddr, *sizeOfBlocks)
	bytes, _ = encryptBytes(bytes, key, iv)
	setBytes(sr, bytes, *startAddr)
	fmt.Println(sr.String())
	bytes, _ = decryptBytes(bytes, key, iv)
	setBytes(sr, bytes, *startAddr)
	fmt.Println(sr.String())
}

func encryptBytes(plainText []byte, key, iv []byte) ([]byte, error) {
	// check length of plainText
	if len(plainText) < aes.BlockSize {
		return []byte{}, fmt.Errorf("ciphertext too short")
	}
	if len(plainText)%aes.BlockSize != 0 {
		return []byte{}, fmt.Errorf("ciphertext is not multiple of the block size")
	}

	// encryption
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	copy(cipherText[:aes.BlockSize], iv)
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	return cipherText[aes.BlockSize:], nil
}

func decryptBytes(cipherText []byte, key, iv []byte) ([]byte, error) {
	// check length of plainText
	if len(cipherText) < aes.BlockSize {
		return []byte{}, fmt.Errorf("ciphertext too short")
	}
	if len(cipherText)%aes.BlockSize != 0 {
		return []byte{}, fmt.Errorf("ciphertext is not multiple of the block size")
	}

	// decryption
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}
	plainText := make([]byte, len(cipherText))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plainText, cipherText)

	return plainText, nil
}

func getBytes(s *srec.Srec, start, size uint32) (bytes []byte) {
loop:
	for _, r := range s.Records {
		if r.Srectype == "S1" || r.Srectype == "S2" || r.Srectype == "S3" {
			current := r.Address
			for _, d := range r.Data {
				if current >= start {
					bytes = append(bytes, d)
					if size--; size == 0 {
						break loop
					}
				}
				current++
			}
		}
	}
	return bytes
}

func setBytes(s *srec.Srec, bytes []byte, start uint32) {
	size := len(bytes)
	index := 0
loop:
	for _, r := range s.Records {
		if r.Srectype == "S1" || r.Srectype == "S2" || r.Srectype == "S3" {
			current := r.Address
			for i := 0; i < len(r.Data); i++ {
				if current >= start {
					r.Data[i] = bytes[index]
					if index++; index == size {
						break loop
					}
				}
				current++
			}
		}
	}
}

func bytesToString(bt []byte) string {
	s := ""
	for i, b := range bt {
		if i != 0 && i%16 == 0 {
			s += fmt.Sprint("\n")
		}
		s += fmt.Sprintf("%02X", b)
	}
	s += fmt.Sprint("\n")
	return s
}
