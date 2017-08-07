package srec

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	TypeFieldStrLen   = 2
	LengthFieldStrLen = 2
	CSumFieldStrLen   = 2
)

type Srec struct {
	HeaderRecord  HeaderRecord
	BinaryRecords []BinaryRecord
	FooterRecord  FooterRecord
	StartAddress  uint32
	EndAddress    uint32
	Bytes         []byte
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
	EntryAddr uint32
	Checksum  byte
}

func NewSrec(outs, errs io.Writer) *Srec {
	return &Srec{OutStream: outs, ErrStream: errs}
}

func NewHeaderRecord() *HeaderRecord {
	return &HeaderRecord{}
}

func NewBianryRecord() *BinaryRecord {
	return &BinaryRecord{}
}

func NewFooterRecord() *FooterRecord {
	return &FooterRecord{}
}

func (srs *Srec) ParseFile(fileReader io.Reader) error {
	rec := new(BinaryRecord)

	scanner := bufio.NewScanner(fileReader)

	for scanner.Scan() {
		line := scanner.Text()
		sl := strings.Split(line, "")

		srectype := strings.Join(sl[:2], "")
		switch {
		case srectype == "S0":
		case (srectype == "S1") || (srectype == "S2") || (srectype == "S3"):
			err := rec.getSrecBinaryRecordFields(srectype, sl)
			if err != nil {
				return err
			}
			srs.BinaryRecords = append(srs.BinaryRecords, *rec)
		case (srectype == "S7") || (srectype == "S8") || (srectype == "S9"):
		default:
			// pass S4~6
		}
	}
	return nil
}

func (rec *BinaryRecord) getSrecBinaryRecordFields(srectype string, sl []string) error {
	var err error

	rec.Srectype = srectype
	rec.Length, err = getLengh(sl)
	if err != nil {
		return err
	}
	rec.Address, err = getAddress(srectype, sl)
	if err != nil {
		return err
	}
	rec.Data, err = getData(srectype, sl)
	if err != nil {
		return err
	}
	rec.Checksum, err = getChecksum(srectype, sl)
	if err != nil {
		return err
	}
	return nil
}

func getAddrLenAsStr(srectype string) (int, error) {
	switch srectype {
	case "S1":
		return 4, nil
	case "S2":
		return 6, nil
	case "S3":
		return 8, nil
	default:
		return 0, fmt.Errorf("%s is not srectype", srectype)
	}
}

func getDataLenAsStr(sl []string) (int, error) {
	len, err := strconv.ParseUint(strings.Join(sl[2:4], ""), 16, 32)
	return int(len * 2), err
}

func getAddress(srectype string, sl []string) (uint32, error) {
	addrLenAsStr, err := getAddrLenAsStr(srectype)
	if err != nil {
		return 0, err
	}
	addr, err := strconv.ParseUint(strings.Join(sl[4:4+addrLenAsStr], ""), 16, 32)
	if err != nil {
		return 0, err
	}
	return uint32(addr), err
}

func getLengh(sl []string) (uint32, error) {
	len, err := strconv.ParseUint(strings.Join(sl[2:4], ""), 16, 32)
	if err != nil {
		return 0, err
	}
	return uint32(len), err
}

func getData(srectype string, sl []string) ([]byte, error) {
	addrLenAsStr, err := getAddrLenAsStr(srectype)
	if err != nil {
		return []byte{}, err
	}
	dataLenAsStr, err := getDataLenAsStr(sl)
	if err != nil {
		return []byte{}, err
	}

	data := make([]byte, 0)
	DataIndexSt := TypeFieldStrLen + LengthFieldStrLen + addrLenAsStr
	DataIndexEd := (TypeFieldStrLen + LengthFieldStrLen) + (dataLenAsStr - CSumFieldStrLen)
	for i := DataIndexSt; i < DataIndexEd; i += 2 {
		b, err := strconv.ParseUint(strings.Join(sl[i:i+2], ""), 16, 32)
		if err != nil {
			return []byte{}, err
		}
		data = append(data, byte(b))
	}
	return data, nil
}

func getChecksum(srectype string, sl []string) (byte, error) {
	dataLenAsStr, err := getDataLenAsStr(sl)
	if err != nil {
		return 0, err
	}

	CSumIndexSt := TypeFieldStrLen + LengthFieldStrLen + dataLenAsStr - CSumFieldStrLen
	CSumIndexEd := TypeFieldStrLen + LengthFieldStrLen + dataLenAsStr
	csum, err := strconv.ParseUint(strings.Join(sl[CSumIndexSt:CSumIndexEd], ""), 16, 32)
	if err != nil {
		return 0, err
	}
	return byte(csum), nil
}

func (sr *Srec) GetBytes() ([]byte, error) {
	sr.StartAddress = getStartAddr(sr)
	sr.EndAddress = getEndAddr(sr)
	LastRecordDatalen := getLastRecordDataLen(sr)
	err := sr.makePaddedBytes(sr.StartAddress, sr.EndAddress, LastRecordDatalen)
	if err != nil {
		return sr.Bytes, err
	}
	return sr.Bytes, err
}

func getStartAddr(sr *Srec) uint32 {
	return sr.BinaryRecords[0].Address
}

func getEndAddr(sr *Srec) uint32 {
	return sr.BinaryRecords[len(sr.BinaryRecords)-1].Address
}

func getLastRecordDataLen(sr *Srec) uint32 {
	len := len(sr.BinaryRecords[len(sr.BinaryRecords)-1].Data)
	return uint32(len)
}

func (sr *Srec) makePaddedBytes(startAddr uint32, endAddr uint32, lastRecordDataLen uint32) error {
	size := (endAddr - startAddr) + lastRecordDataLen
	for i := 0; i < int(size); i++ {
		sr.Bytes = append(sr.Bytes, 0xFF)
	}

	ofst := int(startAddr)
	for _, brcs := range sr.BinaryRecords {
		for i := 0; i < len(brcs.Data); i++ {
			if (brcs.Address < sr.StartAddress) || (brcs.Address > sr.EndAddress) {
				return fmt.Errorf("data address 0x%08X is out of srec range", brcs.Address)
			}
			sr.Bytes[(int(brcs.Address)-ofst)+i] = brcs.Data[i]
		}
	}
	return nil
}
