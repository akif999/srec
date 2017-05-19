package main

import (
	"bufio"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
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

func main() {
	srecs := new(Srecs)

	kingpin.Parse()

	srecs.ParseFile(filename)
	srecs.PrintOnlyData()
	srecs.WriteBinaryToFile()
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
		switch srectype {
		case "S1":
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
			srec.data = data
			srec.checksum = byte(csum)
			sr.records = append(sr.records, *srec)
		case "S2":
		case "S3":
		default:
		}

	}
}

func (sr *Srecs) PrintOnlyData() {
	for _, r := range sr.records {
		for _, b := range r.data {
			fmt.Printf("%02X", b)
		}
		fmt.Println()
	}
}

func (sr *Srecs) WriteBinaryToFile() {
	writeFile, _ := os.OpenFile(*filename+".bin", os.O_WRONLY|os.O_CREATE, 0600)
	writer := bufio.NewWriter(writeFile)
	for _, r := range sr.records {
		writer.Write(r.data)
		writer.Flush()
	}
}
