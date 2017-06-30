package srec

import (
	"reflect"
	"strings"
	"testing"
)

const ()

var ()

func TestGetSrecBinaryFields(t *testing.T) {
	t1Input := strings.Split("S11300E00000010000000100000001000000010008", "")
	t1Want := &BinaryRecord{
		Srectype: "S1",
		Length:   0x13,
		Address:  0x00E0,
		Data: []byte{0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00,
			0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00},
		Checksum: 0x08,
	}
	t1Got := new(BinaryRecord)

	t2Input := strings.Split("S10700F00000010007", "")
	t2Want := &BinaryRecord{
		Srectype: "S1",
		Length:   0x07,
		Address:  0x00F0,
		Data:     []byte{0x00, 0x00, 0x01, 0x00},
		Checksum: 0x07,
	}
	t2Got := new(BinaryRecord)

	t1Got.getSrecBinaryFields(strings.Join(t1Input[:2], ""), t1Input)
	t2Got.getSrecBinaryFields(strings.Join(t2Input[:2], ""), t2Input)
	if reflect.DeepEqual(t1Want, t1Got) != true {
		t.Errorf(" got : %v\n         want : %v", t1Want, t1Got)
	}
	if reflect.DeepEqual(t2Want, t2Got) != true {
		t.Errorf(" got : %v\n         want : %v", t2Want, t2Got)
	}
}
