# srec
A library of Motolola Hex(S Record, *.mot, *mhx) file utilities

## Usage
Now, you can use two interfaces,`GetBytes()` and `SetBytes()`

```go
srec := srec.NewSrec()

srec.ParseFile(fp)
srec.GetBytes()    // it returns all of srec data bytes. you can do this by address in next step
srec.SetBytes(0x00123456, []byte{0x12, 34, 56, 78})
```

## Example
```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AKIF999/srec"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	filename = kingpin.Arg("Filename", "Srec filename").ExistingFile()
	setAddr  = kingpin.Arg("SetAddress", "Address of setting Bytes").Uint32()
)

func main() {
	sr := srec.NewSrec()

	kingpin.Parse()

	fp, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	err = sr.ParseFile(fp)
	if err != nil {
		log.Fatal(err)
	}
	bt := sr.GetBytes()
	for i, b := range bt {
		if i != 0 && i%16 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02X", b)
	}
	fmt.Print("\n\n")

	err = sr.SetBytes(*setAddr, []byte{0x12, 0x34, 0x56, 0x78})
	if err != nil {
		log.Fatal(err)
	}
	bt = sr.GetBytes()
	for i, b := range bt {
		if i != 0 && i%16 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02X", b)
	}
}
```

```
# input file
S00E000074657374202020206D6F7461
S2140D0000CD060D00D6060D00DD060D00E4060D002E
S2140D0010EB060D00F2060D00F9060D0000070D00AB
S2080D002024000D0099
S2140D0024EB600005EB708006C7030A00C7010400E9
S2140D0034B70500C70F0800C7100600B70700B70AB4
S9030000FC
```

```
$./sample input_file 0x0D0020
CD060D00D6060D00DD060D00E4060D00
EB060D00F2060D00F9060D0000070D00
24000D00EB600005EB708006C7030A00
C7010400B70500C70F0800C7100600B7
0700B70A

CD060D00D6060D00DD060D00E4060D00
EB060D00F2060D00F9060D0000070D00
12345678EB600005EB708006C7030A00
C7010400B70500C70F0800C7100600B7
0700B70A
```

## Installation
`go get github.com/AKIF999/srec`

## License
MIT

## Author
Akifumi Kitabatake(user name is "AKIF" or "AKIF999")
