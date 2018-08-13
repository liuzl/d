package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb"
	"github.com/liuzl/d"
	"github.com/liuzl/goutil"
)

var (
	src = flag.String("src", "dict.txt", "dictionary source file")
	dst = flag.String("dst", "dict", "destination dir")
	tag = flag.String("tag", "tag", "word type of this dict")
)

func main() {
	flag.Parse()
	count, err := goutil.FileLineCount(*src)
	if err != nil {
		log.Fatal(err)
	}
	if count <= 0 {
		log.Fatal(fmt.Errorf("empty src file: %s", *src))
	}
	dict, err := d.Load(*dst)
	if err != nil {
		log.Fatal(err)
	}
	fd, _ := os.Open(*src)
	defer fd.Close()
	br := bufio.NewReader(fd)
	bar := pb.StartNew(count)
	value := map[string]interface{}{*tag: nil}
	for {
		word, c := br.ReadString('\n')
		if c == io.EOF {
			break
		}
		if err = dict.Update(word, value); err != nil {
			log.Fatal(err)
		}
		bar.Increment()
	}
	bar.FinishPrint("done!")
	if err = dict.Save(); err != nil {
		log.Fatal(err)
	}
}
