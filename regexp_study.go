package main

import (
	"bufio"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
	"regexp"
)

const ()

var (
	filename = kingpin.Arg("filename", "srec file").ExistingFile()
)

func main() {
	kingpin.Parse()
	re := regexp.MustCompile("S1")

	fp, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		s := scanner.Text()
		if re.MatchString(s) {
			fmt.Println(s)
		}
	}
}
