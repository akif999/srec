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

type Srec struct {
	records []Record
}

type Record struct {
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
	srec := new(Srec)

	kingpin.Parse()

	srec.ParseFile(filename)
	srec.PrintOnlyData()
	srec.WriteBinaryToFile()
}

func (srs *Srec) ParseFile(file *string) {
	rec := new(Record)

	fp, err := os.Open(*file)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	for scanner.Scan() {
		line := scanner.Text()
		sl := strings.Split(line, "")

		srectype := strings.Join(sl[:2], "")
		rec.getSrecFields(srectype, sl)
		if srectype == "S1" {
			srs.records = append(srs.records, *rec)
		}
	}
}

func (rec *Record) getSrecFields(srectype string, sl []string) {
	var len uint64
	var addr uint64
	var data []byte
	var csum uint64

	switch srectype {
	case "S1":
		len, _ = strconv.ParseUint(strings.Join(sl[2:4], ""), 16, 32)
		addr, _ = strconv.ParseUint(strings.Join(sl[4:8], ""), 16, 32)
		data = make([]byte, 0)
		for i := 0; i < (4 + (int(len) * 2) - 2); i += 2 {
			if i >= 8 {
				b, _ := strconv.ParseUint(strings.Join(sl[i:i+2], ""), 16, 32)
				data = append(data, byte(b))
			}
		}
		csum, _ = strconv.ParseUint(strings.Join(sl[4+(int(len)*2)-2:(4+(int(len)*2)-2)+2], ""), 16, 32)
	case "S2":
	case "S3":
	default:
		return
	}
	rec.srectype = srectype
	rec.length = uint32(len)
	rec.address = uint32(addr)
	rec.data = data
	rec.checksum = byte(csum)
}

func (sr *Srec) PrintOnlyData() {
	for _, r := range sr.records {
		for _, b := range r.data {
			fmt.Printf("%02X", b)
		}
		fmt.Println()
	}
}

func (sr *Srec) WriteBinaryToFile() {
	writeFile, _ := os.OpenFile(*filename+".bin", os.O_WRONLY|os.O_CREATE, 0600)
	writer := bufio.NewWriter(writeFile)
	for _, r := range sr.records {
		writer.Write(r.data)
		writer.Flush()
	}
}
