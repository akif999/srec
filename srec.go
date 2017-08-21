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
	length   uint8
	data     []byte
	checksum uint8
}

type dataRecord struct {
	srectype string
	length   uint8
	address  uint32
	data     []byte
	checksum uint8

	isBlank bool
}

type footerRecord struct {
	srectype  string
	length    uint8
	entryAddr uint32
	checksum  uint8
}

func NewSrec() *Srec {
	return &Srec{}
}

func newHeaderRecord() *headerRecord {
	return &headerRecord{}
}

func newDataRecord() *dataRecord {
	return &dataRecord{}
}

func newFooterRecord() *footerRecord {
	return &footerRecord{}
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

func calcChecksum(srectype string, len uint8, addr uint32, data []byte) (uint8, error) {
	fs := ""
	switch srectype {
	case "S1":
		fs = fmt.Sprintf("%02X%04X", len, addr)
	case "S2":
		fs = fmt.Sprintf("%02X%06X", len, addr)
	case "S3":
		fs = fmt.Sprintf("%02X%08X", len, addr)
	default:
	}
	for _, b := range data {
		fs += fmt.Sprintf("%02X", b)
	}
	bs, err := strToBytes(fs)
	if err != nil {
		return 0, err
	}
	sum := byte(0)
	for _, b := range bs {
		sum += b
	}
	return uint8(^sum), nil
}

func strToBytes(src string) ([]byte, error) {
	sp := strings.Split(src, "")
	bs := []byte{}
	for i := 0; i < len(sp); i += 2 {
		b, err := strconv.ParseUint(sp[i]+sp[i+1], 16, 32)
		if err != nil {
			return []byte{}, err
		}
		bs = append(bs, byte(b))
	}
	return bs, nil
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
			rec := newDataRecord()
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

	err = srs.addBlankRecord()
	if err != nil {
		return err
	}

	err = srs.makePaddedBytes(srs.startAddress, srs.endAddress, LastRecordDatalen)
	if err != nil {
		return err
	}
	return nil
}

func (rec *headerRecord) getHeaderRecordFields(sl []string) error {
	var err error

	srectype := "S0"
	rec.length, err = getLength(sl)
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
	rec.length, err = getLength(sl)
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
	rec.isBlank = false
	return nil
}

func (rec *footerRecord) getFooterRecordFields(srectype string, sl []string) error {
	var err error

	rec.srectype = srectype
	rec.length, err = getLength(sl)
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

func getLength(sl []string) (uint8, error) {
	len, err := strconv.ParseUint(strings.Join(sl[2:4], ""), 16, 32)
	if err != nil {
		return 0, err
	}
	return uint8(len), nil
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

func (sr *Srec) isDataRecordExists() error {
	if len(sr.dataRecords) == 0 {
		return fmt.Errorf("byte data is empty. srec file doesn't have S1~3 records.")
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

func (sr *Srec) addBlankRecord() error {
	size := len(sr.dataRecords)
	for i := 0; i < size; i++ {
		if i == size-1 {
			break
		}
		if sr.dataRecords[i].isBlank {
			continue
		}
		cr := sr.dataRecords[i]
		nr := sr.dataRecords[i+1]
		blankSize := (nr.address) - (cr.address + uint32(len(cr.data)))
		if blankSize != 0 {
			addr := cr.address + uint32(len(cr.data))
			for blankSize != 0 {
				dataSize := uint32(16)
				if blankSize < 16 {
					dataSize = blankSize
				}
				r, err := makeBlankRecord(cr.srectype, addr, dataSize)
				if err != nil {
					return err
				}
				sr.dataRecords = append(sr.dataRecords[:i+2], sr.dataRecords[i+1:]...)
				sr.dataRecords[i+1] = r
				blankSize -= dataSize
				addr += dataSize
				i++
			}
			size = len(sr.dataRecords)
		}
	}
	return nil
}

func makeBlankRecord(srectype string, addr uint32, dataSize uint32) (*dataRecord, error) {
	r := newDataRecord()
	r.srectype = srectype
	l, err := getAddrLenAsStr(srectype)
	if err != nil {
		return nil, err
	}
	r.length = uint8(l) + uint8(dataSize*2)
	r.address = addr
	for i := 0; i < int(dataSize); i++ {
		r.data = append(r.data, 0xFF)
	}
	r.checksum, err = calcChecksum(r.srectype, r.length, r.address, r.data)
	if err != nil {
		return nil, err
	}
	r.isBlank = true
	return r, nil
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
