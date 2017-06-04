package main

import (
	"./srec"
	"gopkg.in/alecthomas/kingpin.v2"
)

const ()

var (
	filename = kingpin.Arg("filename", "srec file").ExistingFile()
)

func main() {
	srec := new(srec.Srec)

	kingpin.Parse()

	srec.ParseFile(filename)
	srec.PrintOnlyData()
	srec.WriteBinaryToFile(filename)
}
