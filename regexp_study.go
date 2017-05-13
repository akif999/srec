package main

import (
	"bufio"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	// "regexp"
)

const ()

var (
	file = kingpin.Arg("file", "srec file").ExistingFile()
)

func main() {
	kingpin.Parse()
	// re := regexp.MustCompile('.')

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
