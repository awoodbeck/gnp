package main

import (
	"crypto/sha512"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s file...\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	for _, file := range flag.Args() {
		fmt.Printf("\n%s =>\n%s\n", file, checksum(file))
	}
}

func checksum(file string) string {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%x", sha512.Sum512_256(b))
}
