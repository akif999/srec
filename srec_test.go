package srec

import (
	"reflect"
	"strings"
	"testing"
)

const ()

var ()

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
