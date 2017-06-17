package srec

import (
	"reflect"
	"strings"
	"testing"
)

const ()

var ()

func TestGetSrecBinaryFields(t *testing.T) {
	input := strings.Split("S11300E00000010000000100000001000000010008", "")
	want := &binaryRecord{
		srectype: "S1",
		length:   0x13,
		address:  0x00E0,
		data: []byte{0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00,
			0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00},
		checksum: 0x08,
	}
	got := new(binaryRecord)

	got.getSrecBinaryFields(strings.Join(input[:2], ""), input)
	if reflect.DeepEqual(want, got) != true {
		t.Errorf(" got : %v\n         want : %v", want, got)
	}
}
