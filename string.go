package srec

import "fmt"

// String return the contents of the entire Srecord.
func (sr *Srec) String() string {
	fs := ""
	for _, r := range sr.Records {
		fs += r.String() + "\n"
	}
	return fs
}

// String returns the contents of the Record. If the Record has wrong srectype,
// it returns "<invalid srectype>"
func (r *Record) String() string {
	fs := ""
	switch r.Srectype {
	case "S0", "S1", "S9":
		fs = fmt.Sprintf("%s%02X%04X", r.Srectype, r.Length, r.Address)
	case "S2", "S8":
		fs = fmt.Sprintf("%s%02X%06X", r.Srectype, r.Length, r.Address)
	case "S3", "S7":
		fs = fmt.Sprintf("%s%02X%08X", r.Srectype, r.Length, r.Address)
	case "S5":
		// Since S5 has no address part, it does not format it
		fs = fmt.Sprintf("%s%02X", r.Srectype, r.Length)
	default:
		return "<invalid srectype>"
	}
	for _, b := range r.Data {
		fs += fmt.Sprintf("%02X", b)
	}
	fs += fmt.Sprintf("%02X", r.Checksum)
	return fs
}
