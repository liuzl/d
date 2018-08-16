package main

import (
	"encoding/csv"
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
	src = flag.String("src", "dict.csv", "dictionary source csv file,")
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

	r := csv.NewReader(fd)
	r.FieldsPerRecord = -1

	var errors []error
	bar := pb.StartNew(count)
	for {
		if record, err := r.Read(); err == io.EOF {
			break
		} else if err != nil {
			errors = append(errors, err)
		} else {
			value, err := dict.Get(record[0])
			if err == nil {
				if _, has := value[*tag]; !has {
					value[*tag] = nil
				}
			} else if err.Error() == "leveldb: not found" {
				value = map[string]interface{}{*tag: nil}
			} else {
				log.Fatal(err)
			}
			if len(record) >= 2 && record[1] != record[0] {
				if value[*tag] == nil {
					value[*tag] = record[1]
				} else {
					switch value[*tag].(type) {
					case string:
						former := value[*tag].(string)
						if former != record[1] {
							value[*tag] = []string{former, record[1]}
						}
					case []string:
						former := value[*tag].([]string)
						dup := false
						for _, v := range former {
							if v == record[1] {
								dup = true
								break
							}
						}
						if !dup {
							value[*tag] = append(former, record[1])
						}
					default:
						log.Println("ERROR type")
					}
				}
			}
			if err = dict.Update(record[0], value); err != nil {
				log.Fatal(err)
			}
		}
		bar.Increment()
	}
	bar.FinishPrint("done!")
	if len(errors) > 0 {
		log.Printf("%d errors:\n", len(errors))
		for _, err := range errors {
			log.Println(err)
		}
	}
	if err = dict.Save(); err != nil {
		log.Fatal(err)
	}
}
