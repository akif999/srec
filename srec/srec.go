package srec

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const ()

type Srec struct {
	headerRecord  headerRecord
	binaryRecords []binaryRecord
	footerRecord  footerRecord
	outStream     io.Writer
	errStream     io.Writer
}

type headerRecord struct {
	length   uint32
	data     []byte
	checksum byte
}

type binaryRecord struct {
	srectype string
	length   uint32
	address  uint32
	data     []byte
	checksum byte
}

type footerRecord struct {
	srectype  string
	startaddr uint32
	checksum  byte
}

var ()

func NewSrec(outs, errs io.Writer) *Srec {
	return &Srec{outStream: outs, errStream: errs}
}

func (srs *Srec) ParseFile(file *string) {
	rec := new(binaryRecord)

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
		switch {
		case srectype == "S0":
		case srectype == "S1":
			rec.getSrecBinaryFields(srectype, sl)
			srs.binaryRecords = append(srs.binaryRecords, *rec)
		case srectype == "S2":
		case srectype == "S3":
		case srectype == "S4":
			// S4 is reserved
		case srectype == "S5":
			// S5 is reserved
		case srectype == "S6":
			// S6 is reserved
		case srectype == "S7":
		case srectype == "S8":
		case srectype == "S9":
		}
	}
}

func (rec *binaryRecord) getSrecBinaryFields(srectype string, sl []string) {
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

func Padding() {
}

func (sr *Srec) PrintOnlyData() {
	for _, r := range sr.binaryRecords {
		for _, b := range r.data {
			fmt.Fprintf(sr.outStream, "%02X", b)
		}
		fmt.Println()
	}
}

func (sr *Srec) WriteBinaryToFile(filename *string) {
	writeFile, _ := os.OpenFile(*filename+".bin", os.O_WRONLY|os.O_CREATE, 0600)
	writer := bufio.NewWriter(writeFile)
	for _, r := range sr.binaryRecords {
		writer.Write(r.data)
		writer.Flush()
	}
}
