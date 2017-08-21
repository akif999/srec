package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AKIF999/srec"
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

	err = sr.ParseFile(fp)
	if err != nil {
		log.Fatal(err)
	}

	bt := sr.Bytes()
	printBytes(bt)
	fmt.Print("\n")

	b := []byte{}
	for i := 0; i < 32; i++ {
		b = append(b, 0x88)
	}
	err = sr.SetBytes(*setAddr, b)
	if err != nil {
		log.Fatal(err)
	}
	bt, err = sr.BytesInPart(*getAddr, *getSize)
	if err != nil {
		log.Fatal(err)
	}
	printBytes(bt)
	fmt.Print("\n")

	sr.UpdateInPart(0x00E0, 0x00FF)
	fs := sr.Format()
	fmt.Print(fs)
}

func printBytes(bt []byte) {
	for i, b := range bt {
		if i != 0 && i%16 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02X", b)
	}
	fmt.Print("\n")
}
