package srec

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestNewSrec(t *testing.T) {
	got := NewSrec()
	want := &Srec{}

	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf(" got : %02X\n         want : %02X", got, want)
	}
}

func TestNewRecord(t *testing.T) {
	got := newRecord()
	want := &Record{}

	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf(" got : %02X\n         want : %02X", got, want)
	}
}

func TestGetAddrLen(t *testing.T) {
	type tp struct {
		srectype string
		want     int
		got      int
	}

	tests := []tp{
		tp{
			srectype: "S0",
			want:     2,
		},
		tp{
			srectype: "S1",
			want:     2,
		},
		tp{
			srectype: "S2",
			want:     3,
		},
		tp{
			srectype: "S3",
			want:     4,
		},
		tp{
			srectype: "S7",
			want:     4,
		},
		tp{
			srectype: "S8",
			want:     3,
		},
		tp{
			srectype: "S9",
			want:     2,
		},
	}
	errTest := tp{
		srectype: "S4",
		want:     0,
	}
	for _, tp := range tests {
		tp.got, _ = getAddrLen(tp.srectype)
		if tp.got != tp.want {
			t.Errorf(" got : %02X\n         want : %02X", tp.got, tp.want)
		}
	}
	var err error
	errTest.got, err = getAddrLen(errTest.srectype)
	if err == nil || errTest.got != errTest.want {
		t.Errorf("this method could't detect error.")
	}
}

func TestCalcChecksum(t *testing.T) {
	type tp struct {
		srectype string
		len      uint8
		addr     uint32
		data     []byte
		want     uint8
		got      uint8
	}
	tests := []tp{
		tp{
			srectype: "S1",
			len:      0x13,
			addr:     0x0100,
			data: []byte{
				0x7A, 0x07, 0x00, 0x0F, 0xFF, 0x0E, 0x7A, 0x00,
				0x00, 0x00, 0x01, 0x62, 0x7A, 0x01, 0x00, 0x0F,
			},
			want: 0xE7,
		},
		tp{
			srectype: "S2",
			len:      0x14,
			addr:     0x0D0010,
			data: []byte{
				0xEB, 0x06, 0x0D, 0x00, 0xF2, 0x06, 0x0D, 0x00,
				0xF9, 0x06, 0x0D, 0x00, 0x00, 0x07, 0x0D, 0x00,
			},
			want: 0xAB,
		},
		tp{
			srectype: "S3",
			len:      0x15,
			addr:     0xCAFE0120,
			data: []byte{
				0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55,
				0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55,
			},
			want: 0x09,
		},
		tp{
			srectype: "S2",
			len:      0x08,
			addr:     0x0D0020,
			data: []byte{
				0x24, 0x00, 0x0D, 0x00,
			},
			want: 0x99,
		},
		tp{
			srectype: "S5",
			len:      0x03,
			addr:     0x0000,
			data:     []byte{0x00, 0x11},
			want:     0xEB,
		},
		tp{
			srectype: "S7",
			len:      0x05,
			addr:     0x00000000,
			data:     []byte{},
			want:     0xFA,
		},
		tp{
			srectype: "S9",
			len:      0x03,
			addr:     0x0100,
			data:     []byte{},
			want:     0xFB,
		},
	}
	for _, tp := range tests {
		tp.got, _ = calcChecksum(tp.srectype, tp.len, tp.addr, tp.data)
		if tp.got != tp.want {
			t.Errorf(" got : %02X\n         want : %02X", tp.got, tp.want)
		}
	}
}

func TestEndAddress(t *testing.T) {
	type tp struct {
		rec  *Record
		want uint32
		got  uint32
	}
	tests := []tp{
		tp{
			rec: &Record{
				Srectype: "S1",
				Length:   0x13,
				Address:  0x00E0,
				Data: []byte{
					0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00,
					0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00,
				},
				Checksum: 0x08,
			},
			want: 0x00EF,
		},
		tp{
			rec: &Record{
				Srectype: "S2",
				Length:   0x14,
				Address:  0x0D0010,
				Data: []byte{
					0xEB, 0x06, 0x0D, 0x00, 0xF2, 0x06, 0x0D, 0x00,
					0xF9, 0x06, 0x0D, 0x00, 0x00, 0x07, 0x0D, 0x00,
				},
				Checksum: 0xAB,
			},
			want: 0x0D001F,
		},
		tp{
			rec: &Record{
				Srectype: "S3",
				Length:   0x15,
				Address:  0xCAFE0130,
				Data: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
				Checksum: 0xF1,
			},
			want: 0xCAFE013F,
		},
	}
	for _, tp := range tests {
		tp.got = tp.rec.EndAddress()
		if tp.got != tp.want {
			t.Errorf("got  : %02X\nwant : %02X", tp.got, tp.want)
		}
	}
}

func TestParse(t *testing.T) {
	type tp struct {
		line string
		want *Record
		got  *Record
	}
	tests := []tp{
		tp{
			line: "S11300E00000010000000100000001000000010008",
			want: &Record{
				Srectype: "S1",
				Length:   0x13,
				Address:  0x00E0,
				Data: []byte{
					0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00,
					0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00,
				},
				Checksum: 0x08,
			},
		},

		tp{
			line: "S10700F00000010007",
			want: &Record{
				Srectype: "S1",
				Length:   0x07,
				Address:  0x00F0,
				Data:     []byte{0x00, 0x00, 0x01, 0x00},
				Checksum: 0x07,
			},
		},
		tp{
			line: "S2140D0010EB060D00F2060D00F9060D0000070D00AB",
			want: &Record{
				Srectype: "S2",
				Length:   0x14,
				Address:  0x0D0010,
				Data: []byte{
					0xEB, 0x06, 0x0D, 0x00, 0xF2, 0x06, 0x0D, 0x00,
					0xF9, 0x06, 0x0D, 0x00, 0x00, 0x07, 0x0D, 0x00,
				},
				Checksum: 0xAB,
			},
		},
		tp{
			line: "S2080D002024000D0099",
			want: &Record{
				Srectype: "S2",
				Length:   0x08,
				Address:  0x0D0020,
				Data:     []byte{0x24, 0x00, 0x0D, 0x00},
				Checksum: 0x99,
			},
		},
		tp{
			line: "S315CAFE013000000000000000000000000000000000F1",
			want: &Record{
				Srectype: "S3",
				Length:   0x15,
				Address:  0xCAFE0130,
				Data: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
				Checksum: 0xF1,
			},
		},
		tp{
			line: "S5030011EB",
			want: &Record{
				Srectype: "S5",
				Length:   0x03,
				Address:  0x0,
				Data:     []byte{0x00, 0x11},
				Checksum: 0xEB,
			},
		},
		tp{
			line: "S70500000000FA",
			want: &Record{
				Srectype: "S7",
				Length:   0x05,
				Address:  0x00000000,
				Data:     []byte{},
				Checksum: 0xFA,
			},
		},
		tp{
			line: "S8041234565F",
			want: &Record{
				Srectype: "S8",
				Length:   0x04,
				Address:  0x123456,
				Data:     []byte{},
				Checksum: 0x5F,
			},
		},
		tp{
			line: "S9030100FB",
			want: &Record{
				Srectype: "S9",
				Length:   0x03,
				Address:  0x0100,
				Data:     []byte{},
				Checksum: 0xFB,
			},
		},
	}
	// tests for getRecordFields()
	for _, tp := range tests {
		tp.got = new(Record)
		tp.got.getRecordFields(tp.line)
		if reflect.DeepEqual(tp.want, tp.got) != true {
			t.Errorf(" got : %v\n         want : %v", tp.want, tp.got)
		}
	}
	// test for Parse()
	lines := ""
	want := NewSrec()
	for _, tp := range tests {
		lines += tp.line + "\n"
		want.Records = append(want.Records, tp.want)
	}
	got := NewSrec()
	got.Parse(strings.NewReader(lines))
	if reflect.DeepEqual(want, got) != true {
		t.Errorf(" got : %v\n         want : %v", want, got)
	}
}

func TestEndAddr(t *testing.T) {
	sample := `S00F00006C65645F746573742E6D6F741E
S113004000000100000001000000010000000100A8
S11300500000010000000100000001000000010098
S10700F00000010007
S5030011EB
S11301007A07000FFF0E7A00000001627A01000FE7
S113015046EA01006F62FFFC0B0201006FE2FFFC44
S105016040E673
S9030100FB`

	type tp struct {
		srec *Srec
		want uint32
		got  uint32
	}
	s := NewSrec()
	s.Parse(strings.NewReader(sample))
	tests := []tp{
		tp{
			srec: s,
			want: 0x0161,
		},
	}
	for _, tp := range tests {
		tp.got = tp.srec.EndAddr()
		if tp.got != tp.want {
			t.Errorf("got  : %02X\nwant : %02X", tp.got, tp.want)
		}
	}
}

func TestMakeRec(t *testing.T) {
	type tp struct {
		srectype string
		addr     uint32
		data     []byte
		want     *Record
		got      *Record
	}
	tests := []tp{
		tp{
			srectype: "S1",
			addr:     0x1234,
			data: []byte{
				0x12, 0x34, 0x56, 0x78,
			},
			want: &Record{
				Srectype: "S1",
				Length:   0x07,
				Address:  0x1234,
				Data: []byte{
					0x12, 0x34, 0x56, 0x78,
				},
				Checksum: 0x9E,
			},
		},
		tp{
			srectype: "S2",
			addr:     0x0D0010,
			data: []byte{
				0xEB, 0x06, 0x0D, 0x00, 0xF2, 0x06, 0x0D, 0x00,
				0xF9, 0x06, 0x0D, 0x00, 0x00, 0x07, 0x0D, 0x00,
			},
			want: &Record{
				Srectype: "S2",
				Length:   0x14,
				Address:  0x0D0010,
				Data: []byte{
					0xEB, 0x06, 0x0D, 0x00, 0xF2, 0x06, 0x0D, 0x00,
					0xF9, 0x06, 0x0D, 0x00, 0x00, 0x07, 0x0D, 0x00,
				},
				Checksum: 0xAB,
			},
		},
		tp{
			srectype: "S3",
			addr:     0xCAFE0130,
			data: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			want: &Record{
				Srectype: "S3",
				Length:   0x15,
				Address:  0xCAFE0130,
				Data: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				Checksum: 0xF1,
			},
		},
		tp{
			srectype: "S5",
			addr:     0x0,
			data: []byte{
				0x00, 0x11,
			},
			want: &Record{
				Srectype: "S5",
				Length:   0x03,
				Address:  0x0,
				Data:     []byte{0x00, 0x11},
				Checksum: 0xEB,
			},
		},
		tp{
			srectype: "S7",
			addr:     0x00000000,
			data:     []byte{},
			want: &Record{
				Srectype: "S7",
				Length:   0x05,
				Address:  0x00000000,
				Data:     []byte{},
				Checksum: 0xFA,
			},
		},
		tp{
			srectype: "S8",
			addr:     0x123456,
			data:     []byte{},
			want: &Record{
				Srectype: "S8",
				Length:   0x04,
				Address:  0x123456,
				Data:     []byte{},
				Checksum: 0x5F,
			},
		},
		tp{
			srectype: "S9",
			addr:     0x0100,
			data:     []byte{},
			want: &Record{
				Srectype: "S9",
				Length:   0x03,
				Address:  0x0100,
				Data:     []byte{},
				Checksum: 0xFB,
			},
		},
		tp{
			srectype: "S1",
			addr:     0xFFFC,
			data: []byte{
				0x12, 0x34, 0x56, 0x78,
			},
			want: &Record{
				Srectype: "S1",
				Length:   0x07,
				Address:  0x00FFFC,
				Data: []byte{
					0x12, 0x34, 0x56, 0x78,
				},
				Checksum: 0xE9,
			},
		},
		tp{
			srectype: "S1",
			addr:     0xFFFD,
			data: []byte{
				0x12, 0x34, 0x56, 0x78,
			},
			want: &Record{
				Srectype: "S2",
				Length:   0x08,
				Address:  0x00FFFD,
				Data: []byte{
					0x12, 0x34, 0x56, 0x78,
				},
				Checksum: 0xE7,
			},
		},
		tp{
			srectype: "S2",
			addr:     0xFFFFFC,
			data: []byte{
				0x12, 0x34, 0x56, 0x78,
			},
			want: &Record{
				Srectype: "S2",
				Length:   0x08,
				Address:  0x00FFFFFC,
				Data: []byte{
					0x12, 0x34, 0x56, 0x78,
				},
				Checksum: 0xE9,
			},
		},
		tp{
			srectype: "S2",
			addr:     0xFFFFFD,
			data: []byte{
				0x12, 0x34, 0x56, 0x78,
			},
			want: &Record{
				Srectype: "S3",
				Length:   0x09,
				Address:  0x00FFFFFD,
				Data: []byte{
					0x12, 0x34, 0x56, 0x78,
				},
				Checksum: 0xE7,
			},
		},
	}
	for _, tp := range tests {
		tp.got, _ = MakeRec(tp.srectype, tp.addr, tp.data)
		if reflect.DeepEqual(tp.got, tp.want) != true {
			t.Errorf("got : %v\nwant : %v", tp.got, tp.want)
		}
	}
}

func TestScanLinesCustom(t *testing.T) {
	type tp struct {
		input string
		want  string
		got   string
	}
	tests := []tp{
		// 1~4 : no lb or one at end
		tp{
			input: "abcdefghijkl",
			want:  "abcdefghijkl\n",
		},
		tp{
			input: "abcdefghijkl\r\n",
			want:  "abcdefghijkl\n",
		},
		tp{
			input: "abcdefghijkl\n",
			want:  "abcdefghijkl\n",
		},
		tp{
			input: "abcdefghijkl\r",
			want:  "abcdefghijkl\n",
		},
		// 5~7 : top
		tp{
			input: "\r\nabcdefghijkl",
			want:  "\nabcdefghijkl\n",
		},
		tp{
			input: "\nabcdefghijkl",
			want:  "\nabcdefghijkl\n",
		},
		tp{
			input: "\rabcdefghijkl",
			want:  "\nabcdefghijkl\n",
		},
		// 8~11 : top and buttom
		tp{
			input: "\r\nabcdefghijkl\r\n",
			want:  "\nabcdefghijkl\n",
		},
		tp{
			input: "\nabcdefghijkl\n",
			want:  "\nabcdefghijkl\n",
		},
		tp{
			input: "\rabcdefghijkl\r",
			want:  "\nabcdefghijkl\n",
		},
		tp{
			input: "\r\nabcdefghijkl\n",
			want:  "\nabcdefghijkl\n",
		},
		// 12 : only crlf
		tp{
			input: "abc\r\ndef\r\nghi\r\njkl\r\n",
			want:  "abc\ndef\nghi\njkl\n",
		},
		// 13: only lf
		tp{
			input: "abc\ndef\nghi\njkl\n",
			want:  "abc\ndef\nghi\njkl\n",
		},
		// 14 : only cr
		tp{
			input: "abc\rdef\rghi\rjkl\r",
			want:  "abc\ndef\nghi\njkl\n",
		},
		// 15 :  lf in crlf
		tp{
			input: "abc\r\ndef\nghi\r\njkl\r\n",
			want:  "abc\ndef\nghi\njkl\n",
		},
		// 16 : cr in crlf
		tp{
			input: "abc\r\ndef\rghi\r\njkl\r\n",
			want:  "abc\ndef\nghi\njkl\n",
		},
		// 17 : crlf in lf
		tp{
			input: "abc\ndef\nghi\r\njkl\n",
			want:  "abc\ndef\nghi\njkl\n",
		},
		// 18 : cr in lf
		tp{
			input: "abc\ndef\nghi\rjkl\n",
			want:  "abc\ndef\nghi\njkl\n",
		},
		// 19 : crlf in cr
		tp{
			input: "abc\rdef\rghi\r\njkl\r",
			want:  "abc\ndef\nghi\njkl\n",
		},
		// 20 : lf in cr
		tp{
			input: "abc\rdef\rghi\njkl\r",
			want:  "abc\ndef\nghi\njkl\n",
		},
		// 21 : crlf duplicate
		tp{
			input: "abc\r\ndef\r\nghi\r\n\r\njkl\r\n",
			want:  "abc\ndef\nghi\n\njkl\n",
		},
		// 22 : lf duplicate
		tp{
			input: "abc\ndef\nghi\n\njkl\n",
			want:  "abc\ndef\nghi\n\njkl\n",
		},
		// 23 : cr duplicate
		tp{
			input: "abc\rdef\rghi\r\rjkl\r",
			want:  "abc\ndef\nghi\n\njkl\n",
		},

		// 24 : cr duplicate in crlf
		tp{
			input: "abc\r\ndef\r\nghi\r\r\njkl\r\n",
			want:  "abc\ndef\nghi\n\njkl\n",
		},
		// 25 : lf duplicate in crlf
		tp{
			input: "abc\r\ndef\r\nghi\r\n\njkl\r\n",
			want:  "abc\ndef\nghi\n\njkl\n",
		},
		// 26 : crlf duplicate in crlf
		tp{
			input: "abc\r\ndef\r\nghi\r\n\r\njkl\r\n",
			want:  "abc\ndef\nghi\n\njkl\n",
		},
	}
	for i, e := range tests {
		e.got = ""
		scanner := bufio.NewScanner(strings.NewReader(e.input))
		scanner.Split(scanLinesCustom)
		for scanner.Scan() {
			e.got += scanner.Text() + "\n"
		}
		if e.got != e.want {
			t.Errorf("case %d got :\n%s\nwant :\n%s\n", i+1, e.got, e.want)
		}
	}
}
