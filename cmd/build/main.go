package main

import (
	"flag"
	"log"

	"github.com/liuzl/d"
)

var (
	src = flag.String("src", "dict.json", "dictionary source file")
	dst = flag.String("dst", "dict", "destination dir")
)

func main() {
	flag.Parse()
	_, err := d.Build(*src, *dst)
	if err != nil {
		log.Fatal(err)
	}
}
