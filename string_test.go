package srec

import (
	"strings"
	"testing"
)

func TestStringer(t *testing.T) {
	tests := []string{
		`S00A000065616C6C4350467E`,
		`S2140740209BD0FC00B3009F8F0002FFFC871E9FA05B`,
		`S804000000FB`,
	}

	for _, test := range tests {
		s := NewSrec()
		s.Parse(strings.NewReader(test))
		if g, e := s.String(), test+"\n"; g != e {
			t.Errorf("got %q, want %q", g, e)
		}
		if g, e := s.Records[0].String(), test; g != e {
			t.Errorf("got %q, want %q", g, e)
		}
	}
}

func TestStringerCombined(t *testing.T) {
	tests := []string{
		`S00A000065616C6C4350467E`,
		`S2140740209BD0FC00B3009F8F0002FFFC871E9FA05B`,
		`S804000000FB`,
	}

	combined := strings.Join(tests, "\n")
	s := NewSrec()
	s.Parse(strings.NewReader(combined))

	if g, e := s.String(), combined+"\n"; g != e {
		t.Errorf("got %q, want %q", g, e)
	}

	for i, r := range s.Records {
		if g, e := r.String(), tests[i]; g != e {
			t.Errorf("got %q, want %q", g, e)
		}
	}
}
