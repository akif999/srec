package srec

import "fmt"

func (sr *Srec) Format() string {
	fs := sr.headerRecord.formatHeader()
	for _, r := range sr.dataRecords {
		fs += r.formatData()
	}
	fs += sr.footerRecord.formatFooter()
	return fs
}

func (r *headerRecord) formatHeader() string {
	fs := fmt.Sprintf("S0%02X0000", r.length)
	for _, b := range r.data {
		fs += fmt.Sprintf("%02X", b)
	}
	fs += fmt.Sprintf("%02X\n", r.checksum)
	return fs
}

func (r *dataRecord) formatData() string {
	if r.isBlank {
		return ""
	}
	fs := ""
	switch r.srectype {
	case "S1":
		fs = fmt.Sprintf("S1%02X%04X", r.length, r.address)
	case "S2":
		fs = fmt.Sprintf("S2%02X%06X", r.length, r.address)
	case "S3":
		fs = fmt.Sprintf("S1%02X%08X", r.length, r.address)
	}
	for _, b := range r.data {
		fs += fmt.Sprintf("%02X", b)
	}
	fs += fmt.Sprintf("%02X\n", r.checksum)
	return fs
}

func (r *footerRecord) formatFooter() string {
	fs := ""
	switch r.srectype {
	case "S7":
		fs = fmt.Sprintf("S7%02X%08X", r.length, r.entryAddr)
	case "S8":
		fs = fmt.Sprintf("S8%02X%06X", r.length, r.entryAddr)
	case "S9":
		fs = fmt.Sprintf("S9%02X%04X", r.length, r.entryAddr)
	}
	fs += fmt.Sprintf("%02X\n", r.checksum)
	return fs
}
