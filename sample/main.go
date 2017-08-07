// this source code is sample of package srec
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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
	bt, err := sr.GetBytes()
	if err != nil {
		log.Fatal(err)
	}
	for i, b := range bt {
		if i != 0 && i%16 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02X", b)
	}
	fmt.Println()

	sr2 := srec.NewSrec(c.outs, c.errs)
	file := strings.NewReader("S113000000000100000001000000010000000100E8")
	err = sr2.ParseFile(file)
	if err != nil {
		log.Fatal(err)
	}
	bt2, err := sr2.GetBytes()
	if err != nil {
		log.Fatal(err)
	}
	for i, b := range bt2 {
		if i%16 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02X", b)
	}
	// PrintOnlyData(srec)
	// WriteBinaryToFile(srec, filename)
}

func PrintOnlyData(sr *srec.Srec) {
	for _, r := range sr.BinaryRecords {
		for _, b := range r.Data {
			fmt.Fprintf(sr.OutStream, "%02X", b)
		}
		fmt.Println()
	}
}

func WriteBinaryToFile(sr *srec.Srec, filename *string) {
	writeFile, _ := os.OpenFile(*filename+".bin", os.O_WRONLY|os.O_CREATE, 0600)
	writer := bufio.NewWriter(writeFile)
	for _, r := range sr.BinaryRecords {
		writer.Write(r.Data)
		writer.Flush()
	}
}
