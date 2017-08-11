// this source code is sample of package srec
package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/AKIF999/srec"
	"gopkg.in/alecthomas/kingpin.v2"
)

const ()

var (
	filename = kingpin.Arg("filename", "srec file").ExistingFile()
)

type cli struct {
	outs io.Writer
	errs io.Writer
}

func main() {
	c := &cli{outs: os.Stdout, errs: os.Stderr}
	sr := srec.NewSrec(c.outs, c.errs)

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
	bt := sr.GetBytes()
	for i, b := range bt {
		if i != 0 && i%16 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02X", b)
	}
	fmt.Print("\n\n")

	err = sr.SetBytes(0x000000E4, []byte{0x12, 0x34, 0x56, 0x78})
	if err != nil {
		log.Fatal(err)
	}
	bt = sr.GetBytes()
	for i, b := range bt {
		if i != 0 && i%16 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02X", b)
	}
}
