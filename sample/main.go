// this source code is sample of package srec
package main

import (
	"bufio"
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
	srec := srec.NewSrec(c.outs, c.errs)

	kingpin.Parse()

	fp, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	srec.ParseFile(fp)
	PrintOnlyData(srec)
	WriteBinaryToFile(srec, filename)
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
