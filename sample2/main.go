package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/akif999/srec"
)

var (
	inputFile    = kingpin.Arg("inputFile", "path of input-file(*.mhx)").ExistingFile()
	outputFile   = kingpin.Arg("outputFile", "name of output-file(*.bin)").String()
	startAddr    = kingpin.Arg("StartAddr", "address of start converting").Uint32()
	sizeOfBlocks = kingpin.Arg("sizeOfBlocks", "Size of blocks of converting").Uint32()
)

func main() {
	kingpin.Parse()

	fpIn, err := os.Open(*inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fpIn.Close()

	sr := srec.NewSrec()
	err = sr.Parse(fpIn)
	if err != nil {
		log.Fatal(err)
	}
	data := getBytes(sr, *startAddr, *sizeOfBlocks)

	fpOut, err := os.Create(*outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fpOut.Close()

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		log.Fatal(err)
	}
	fpOut.Write(buf.Bytes())
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
