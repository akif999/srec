package main

import (
	"bufio"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
	"regexp"
)

const ()

type Srec struct {
	srectype uint
	length   uint32
	address  uint32
	data     []byte
	checksum byte
}

var (
	filename = kingpin.Arg("filename", "srec file").ExistingFile()
)

// TODO グループ化してマッチさせ、columnスプリットでフィールドを取り出す
func main() {
	kingpin.Parse()
	re := regexp.MustCompile("S1")

	fp, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	for scanner.Scan() {
		s := scanner.Text()
		if re.MatchString(s) {
			fmt.Println(s)
		}
	}
}
