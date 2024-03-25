# srec

A library of Motolola Hex(S Record, *.mot, *mhx) file utilities

## Usage

You can use following interfaces :  

```
// NewSrec returns a new Srec object
func NewSrec() *Srec {
// EndAddress returns endaddress of the record
func (rec *Record) EndAddress() uint32 {
// CalcChecksum calculates the checksum value of the record from the information of the arguments
func (rec *Record) CalcChecksum() (uint8, error) {
// CalcChecksum calculates the checksum value of the record from the information of the arguments
// S4, S6 are not handled
func (srs *Srec) Parse(fileReader io.Reader) error {
func (srs *Srec) EndAddr() uint32 {
// MakeRec creates and returns a new Record object from the argument information
func MakeRec(srectype string, addr uint32, data []byte) (*Record, error) {
```

## Example

See following files:

[./sample/main.go](./sample/main.go)  
[./sample2/main.go](./sample2/main.go) 

## Installation

`go get github.com/akif999/srec`

## License

MIT

## Author

Akifumi Kitabatake(a.k.a akif999)
