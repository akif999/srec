package srec

import (
	"reflect"
	"strings"
	"testing"
)

const ()

var ()

func TestCalcChecksum(t *testing.T) {
	type tp struct {
		srectype string
		len      uint8
		addr     uint32
		data     []byte
		want     uint8
		got      uint8
	}
	t1 := new(tp)
	t1.srectype = "S1"
	t1.len = 0x13
	t1.addr = 0x0100
	t1.data = []byte{
		0x7A, 0x07, 0x00, 0x0F, 0xFF, 0x0E, 0x7A, 0x00,
		0x00, 0x00, 0x01, 0x62, 0x7A, 0x01, 0x00, 0x0F,
	}
	t1.want = 0xE7
	t1.got, _ = calcChecksum(t1.srectype, t1.len, t1.addr, t1.data)
	if t1.got != t1.want {
		t.Errorf(" got : %02X\n         want : %02X", t1.got, t1.want)
	}
}

func TestGetDataRecordFields(t *testing.T) {
	t1Input := strings.Split("S11300E00000010000000100000001000000010008", "")
	t1Want := &dataRecord{
		srectype: "S1",
		length:   0x13,
		address:  0x00E0,
		data: []byte{0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00,
			0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00},
		checksum: 0x08,
	}
	t1Got := new(dataRecord)

	t2Input := strings.Split("S10700F00000010007", "")
	t2Want := &dataRecord{
		srectype: "S1",
		length:   0x07,
		address:  0x00F0,
		data:     []byte{0x00, 0x00, 0x01, 0x00},
		checksum: 0x07,
	}
	t2Got := new(dataRecord)

	t1Got.getDataRecordFields(strings.Join(t1Input[:2], ""), t1Input)
	t2Got.getDataRecordFields(strings.Join(t2Input[:2], ""), t2Input)
	if reflect.DeepEqual(t1Want, t1Got) != true {
		t.Errorf(" got : %v\n         want : %v", t1Want, t1Got)
	}
	if reflect.DeepEqual(t2Want, t2Got) != true {
		t.Errorf(" got : %v\n         want : %v", t2Want, t2Got)
	}
}
