package srec

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const ()

type Srec struct {
	BeaderRecord  HeaderRecord
	BinaryRecords []BinaryRecord
	BooterRecord  FooterRecord
	BtartAddress  uint
	EndAddress    uint
	OutStream     io.Writer
	ErrStream     io.Writer
}

type HeaderRecord struct {
	Length   uint32
	Data     []byte
	Checksum byte
}

type BinaryRecord struct {
	Srectype string
	Length   uint32
	Address  uint32
	Data     []byte
	Checksum byte
}

type FooterRecord struct {
	Srectype  string
	Startaddr uint32
	Checksum  byte
}

var ()

func NewSrec(outs, errs io.Writer) *Srec {
	return &Srec{OutStream: outs, ErrStream: errs}
}

func (srs *Srec) ParseFile(file *string) {
	rec := new(BinaryRecord)

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
			srs.binaryRecords = append(srs.BinaryRecords, *rec)
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

func (rec *headerRecord) getSrecHeaderFields(srectype string, sl []string) {
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
	rec.Srectype = srectype
	rec.Length = uint32(len)
	rec.Address = uint32(addr)
	rec.Data = data
	rec.Checksum = byte(csum)
}

func (rec *FooterRecord) getSrecFooterFields(srectype string, sl []string) {
}

func Padding() {
}
