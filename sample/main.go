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
	getAddr  = kingpin.Arg("GetAddress", "Address of getting Bytes").Uint32()
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

	nr, err := srec.MakeRec("S1", 0xFFFC, []byte{0x12, 0x34, 0x56, 0x78})
	if err != nil {
		log.Fatal(err)
	}
	sr.Records = append(sr.Records, nr)

	fmt.Print(sr.String())
}
