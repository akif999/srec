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
	for i, b := range bt {
		if i != 0 && i%16 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02X", b)
	}
	fmt.Print("\n\n")

	err = sr.SetBytes(*setAddr, []byte{0x12, 0x34, 0x56, 0x78})
	if err != nil {
		log.Fatal(err)
	}
	bt = sr.Bytes()
	for i, b := range bt {
		if i != 0 && i%16 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02X", b)
	}
}
