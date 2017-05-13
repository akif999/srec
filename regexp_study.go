package main

import (
	"bufio"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
	// "regexp"
	"strings"
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
	// re := regexp.MustCompile("S1")

	fp, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	for scanner.Scan() {
		line := scanner.Text()
		ss := strings.Split(line, "")
		if (ss[0] == "S") && (ss[1] == "1") {
			fmt.Println(ss)
		}
	}
}

func ParseSrec() {
}

func PrintOnlyData() {
}
