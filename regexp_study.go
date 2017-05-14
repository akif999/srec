package main

import (
	"bufio"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
	// "regexp"
	"strconv"
	"strings"
)

const ()

type Srecs struct {
	records []Srec
}

type Srec struct {
	srectype string
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
	srecs := new(Srecs)

	kingpin.Parse()

	srecs.ParseFile(filename)

}

func (sr *Srecs) ParseFile(file *string) {
	srec := new(Srec)

	fp, err := os.Open(*file)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	for scanner.Scan() {
		line := scanner.Text()
		ss := strings.Split(line, "")

		/* get srectype*/
		srectype := strings.Join(ss[:2], "")
		if srectype != "S1" {
			continue
		}
		/* get fields */
		len, _ := strconv.ParseUint(strings.Join(ss[2:4], ""), 16, 32)
		addr, _ := strconv.ParseUint(strings.Join(ss[4:8], ""), 16, 32)
		data := make([]byte, 0)
		for i := 0; i < (4 + (int(len) * 2) - 2); i += 2 {
			if i >= 8 {
				b, _ := strconv.ParseUint(strings.Join(ss[i:i+2], ""), 16, 32)
				data = append(data, byte(b))
			}
		}
		csum, _ := strconv.ParseUint(strings.Join(ss[4+(int(len)*2)-2:(4+(int(len)*2)-2)+2], ""), 16, 32)

		/* put fields */
		srec.srectype = srectype
		srec.length = uint32(len)
		srec.address = uint32(addr)
		srec.checksum = byte(csum)
		srec.data = data
		sr.records = append(sr.records, *srec)

		fmt.Printf("%s %02X %04X ", srec.srectype, srec.length, srec.address)
		for _, b := range srec.data {
			fmt.Printf("%02X", b)
		}
		fmt.Printf(" %02X", srec.checksum)
		fmt.Println()
	}
}

func PrintOnlyData() {
}
