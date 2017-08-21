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
	headerRecord *headerRecord
	dataRecords  []*dataRecord
	footerRecord *footerRecord
	startAddress uint32
	endAddress   uint32
	dataBytes    []byte
}

type headerRecord struct {
	length   uint32
	data     []byte
	checksum uint8
}

type dataRecord struct {
	srectype string
	length   uint32
	address  uint32
	data     []byte
	checksum uint8
}

type footerRecord struct {
	srectype  string
	length    uint32
	entryAddr uint32
	checksum  uint8
}

func NewSrec() *Srec {
	return &Srec{}
}

func newHeaderRecord() *headerRecord {
	return &headerRecord{}
}

func newBianryRecord() *dataRecord {
	return &dataRecord{}
}

func newFooterRecord() *footerRecord {
	return &footerRecord{}
}

func (srs *Srec) ParseFile(fileReader io.Reader) error {
	scanner := bufio.NewScanner(fileReader)

	for scanner.Scan() {
		splitedLine := strings.Split(scanner.Text(), "")

		srectype := strings.Join(splitedLine[:2], "")
		switch {
		case srectype == "S0":
			rec := newHeaderRecord()
			err := rec.getHeaderRecordFields(splitedLine)
			if err != nil {
				return err
			}
			srs.headerRecord = rec
		case (srectype == "S1") || (srectype == "S2") || (srectype == "S3"):
			rec := newBianryRecord()
			err := rec.getDataRecordFields(srectype, splitedLine)
			if err != nil {
				return err
			}
			srs.dataRecords = append(srs.dataRecords, rec)
		case (srectype == "S7") || (srectype == "S8") || (srectype == "S9"):
			rec := newFooterRecord()
			err := rec.getFooterRecordFields(srectype, splitedLine)
			if err != nil {
				return err
			}
			srs.footerRecord = rec
		default:
			// pass S4~6
		}
	}

	err := srs.isDataRecordExists()
	if err != nil {
		return err
	}
	err = srs.isAddrAcending()
	if err != nil {
		return err
	}

	srs.startAddress = getStartAddr(srs)
	srs.endAddress = getEndAddr(srs)
	LastRecordDatalen := getLastRecordDataLen(srs)

	err = srs.makePaddedBytes(srs.startAddress, srs.endAddress, LastRecordDatalen)
	if err != nil {
		return err
	}
	return nil
}

func (rec *headerRecord) getHeaderRecordFields(sl []string) error {
	var err error

	srectype := "S0"
	rec.length, err = getLengh(sl)
	if err != nil {
		return err
	}
	rec.data, err = getData(srectype, sl)
	if err != nil {
		return err
	}
	rec.checksum, err = getChecksum(srectype, sl)
	if err != nil {
		return err
	}
	return nil
}

func (rec *dataRecord) getDataRecordFields(srectype string, sl []string) error {
	var err error

	rec.srectype = srectype
	rec.length, err = getLengh(sl)
	if err != nil {
		return err
	}
	rec.address, err = getAddress(srectype, sl)
	if err != nil {
		return err
	}
	rec.data, err = getData(srectype, sl)
	if err != nil {
		return err
	}
	rec.checksum, err = getChecksum(srectype, sl)
	if err != nil {
		return err
	}
	return nil
}

func (rec *footerRecord) getFooterRecordFields(srectype string, sl []string) error {
	var err error

	rec.srectype = srectype
	rec.length, err = getLengh(sl)
	if err != nil {
		return err
	}
	rec.entryAddr, err = getAddress(srectype, sl)
	if err != nil {
		return err
	}
	rec.checksum, err = getChecksum(srectype, sl)
	if err != nil {
		return err
	}
	return nil
}

func getLengh(sl []string) (uint32, error) {
	len, err := strconv.ParseUint(strings.Join(sl[2:4], ""), 16, 32)
	if err != nil {
		return 0, err
	}
	return uint32(len), nil
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
	return uint32(addr), nil
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

func getChecksum(srectype string, sl []string) (uint8, error) {
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

func getAddrLenAsStr(srectype string) (int, error) {
	switch srectype {
	case "S0":
		return 4, nil
	case "S1":
		return 4, nil
	case "S2":
		return 6, nil
	case "S3":
		return 8, nil
	case "S7":
		return 8, nil
	case "S8":
		return 6, nil
	case "S9":
		return 4, nil
	default:
		return 0, fmt.Errorf("%s is not srectype.", srectype)
	}
}

func getDataLenAsStr(sl []string) (int, error) {
	len, err := strconv.ParseUint(strings.Join(sl[2:4], ""), 16, 32)
	return int(len * 2), err
}

func (sr *Srec) isDataRecordExists() error {
	if len(sr.dataRecords) == 0 {
		return fmt.Errorf("byte data is empty. call PaeseFile() or maybe srec file has no S1~3 records.")
	}
	return nil
}

func (sr *Srec) isAddrAcending() error {
	var prevAddr uint32
	for i, brec := range sr.dataRecords {
		if i == 0 {
			continue
		}
		if brec.address < prevAddr {
			return fmt.Errorf("Address is not acending order.")
		}
		prevAddr = brec.address
	}
	return nil
}

func getStartAddr(sr *Srec) uint32 {
	return sr.dataRecords[0].address
}

func getEndAddr(sr *Srec) uint32 {
	return sr.dataRecords[len(sr.dataRecords)-1].address
}

func getLastRecordDataLen(sr *Srec) uint32 {
	len := len(sr.dataRecords[len(sr.dataRecords)-1].data)
	return uint32(len)
}

func (sr *Srec) makePaddedBytes(startAddr uint32, endAddr uint32, lastRecordDataLen uint32) error {
	size := (endAddr - startAddr) + lastRecordDataLen
	for i := 0; i < int(size); i++ {
		sr.dataBytes = append(sr.dataBytes, 0xFF)
	}

	ofst := int(startAddr)
	for _, brcs := range sr.dataRecords {
		for i := 0; i < len(brcs.data); i++ {
			if (brcs.address < sr.startAddress) || (brcs.address > sr.endAddress) {
				return fmt.Errorf("data address 0x%08X is out of srec range.", brcs.address)
			}
			sr.dataBytes[(int(brcs.address)-ofst)+i] = brcs.data[i]
		}
	}
	return nil
}
