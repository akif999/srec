package main

import (
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
	srec.PrintOnlyData()
	srec.WriteBinaryToFile(filename)
}
