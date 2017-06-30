package main

import (
	"bufio"
	"fmt"
	"github.com/AKIF999/srec"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
)

const ()

var (
	filename = kingpin.Arg("filename", "srec file").ExistingFile()
)

type cli struct {
	outs io.Writer
	errs io.Writer
}

// このmain.goおよびmain()は、ライブラリのテスト用だが、
// srec packageのテストを追加後は不要になる見込み
func main() {
	c := &cli{outs: os.Stdout, errs: os.Stderr}
	srec := srec.NewSrec(c.outs, c.errs)

	kingpin.Parse()

	srec.ParseFile(filename)
	PrintOnlyData(srec)
	WriteBinaryToFile(srec, filename)
}

func PrintOnlyData(sr *srec.Srec) {
	for _, r := range sr.BinaryRecords {
		for _, b := range r.data {
			fmt.Fprintf(sr.OutStream, "%02X", b)
		}
		fmt.Println()
	}
}

func WriteBinaryToFile(sr *srec.Srec, filename *string) {
	writeFile, _ := os.OpenFile(*filename+".bin", os.O_WRONLY|os.O_CREATE, 0600)
	writer := bufio.NewWriter(writeFile)
	for _, r := range sr.BinaryRecords {
		writer.Write(r.data)
		writer.Flush()
	}
}
