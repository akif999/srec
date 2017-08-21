package srec

import "fmt"

const (
	StatOutOfRng = iota
	StatUpdating
	StatFnUpdate
)

func (sr *Srec) Update() error {
	err := sr.UpdateInPart(sr.startAddress, sr.endAddress+uint32(sr.lastRecordDataLen))
	if err != nil {
		return err
	}
	return nil
}

func (sr *Srec) UpdateInPart(stAddr uint32, edAddr uint32) error {
	if stAddr > edAddr {
		return fmt.Errorf("start address must be smaller than end address.")
	}
	if (stAddr < sr.startAddress) || (edAddr > (sr.endAddress + uint32(sr.lastRecordDataLen))) {
		return fmt.Errorf("start address 0x%08X is out of srec range.", stAddr)
	}

	state := StatOutOfRng
	pos := stAddr - sr.startAddress
	for _, r := range sr.dataRecords {
		if state == StatFnUpdate {
			break
		}
		if (stAddr >= r.address) || (stAddr <= (r.address + uint32(len(r.data)))) {
			if r.isBlank {
				r.isBlank = false
			}
			for i := 0; i < len(r.data); i++ {
				if (r.address + uint32(i)) == stAddr {
					state = StatUpdating
				}
				if state == StatUpdating {
					r.data[i] = sr.dataBytes[pos]
					pos++
				}
				if (r.address + uint32(i)) == edAddr {
					state = StatFnUpdate
				}
			}
			var err error
			r.checksum, err = calcChecksum(r.srectype, r.length, r.address, r.data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
