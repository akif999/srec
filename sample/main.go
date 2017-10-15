package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akif999/srec"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	filename = kingpin.Arg("Filename", "Srec filename").ExistingFile()
	setAddr  = kingpin.Arg("SetAddress", "Address of setting Bytes").Uint32()
	getSize  = kingpin.Arg("GetSize", "Size of getting Bytes").Uint32()
)

func main() {
	sr := srec.NewSrec()

	kingpin.Parse()

	fp, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	err = sr.Parse(fp)
	if err != nil {
		log.Fatal(err)
	}

	nr, err := srec.MakeRec("S1", 0x0170, []byte{0x12, 0x34, 0x56, 0x78})
	if err != nil {
		log.Fatal(err)
	}
	sr.Records = append(sr.Records, nr)

	fmt.Println(sr.String())
	fmt.Print(bytesToString(getBytes(sr)))
}

func getBytes(s *srec.Srec) []byte {
	bytes := []byte{}

	size := int32(getEndAddr(s) - getStartAddr(s))
	for size >= 0 {
		bytes = append(bytes, 0xFF)
		size--
	}

	offset := getStartAddr(s)
	for _, r := range s.Records {
		addr := r.Address - offset
		for _, b := range r.Data {
			bytes[addr] = b
			addr++
		}
	}
	return bytes
}

func getStartAddr(s *srec.Srec) uint32 {
	for _, r := range s.Records {
		if r.Srectype == "S1" || r.Srectype == "S2" || r.Srectype == "S3" {
			return r.Address
		}
	}
	return 0x00000000
}

func getEndAddr(s *srec.Srec) uint32 {
	return s.EndAddr()
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
