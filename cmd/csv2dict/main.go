package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb"
	"github.com/golang/glog"
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
		glog.Fatal(err)
	}
	if count <= 0 {
		glog.Fatal(fmt.Errorf("empty src file: %s", *src))
	}
	dict, err := d.Load(*dst)
	if err != nil {
		glog.Fatal(err)
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
			save := false
			value, err := dict.Get(record[0])
			if err == nil {
				if _, has := value[*tag]; !has {
					value[*tag] = nil
					save = true
				}
			} else if err.Error() == "leveldb: not found" {
				value = map[string]interface{}{*tag: nil}
				save = true
			} else {
				glog.Fatal(err)
			}
			if len(record) >= 2 && record[1] != record[0] {
				if value[*tag] == nil {
					value[*tag] = []string{record[1]}
				} else {
					switch value[*tag].(type) {
					case []string:
						values := value[*tag].([]string)
						dup := false
						for _, v := range values {
							if v == record[1] {
								dup = true
								break
							}
						}
						if !dup {
							value[*tag] = append(values, record[1])
							save = true
						}
					default:
						glog.Error("ERROR type")
					}
				}
			}
			if save {
				if err = dict.Update(record[0], value); err != nil {
					glog.Fatalf("%+v, %+v", record, err)
				}
			}
		}
		bar.Increment()
	}
	bar.FinishPrint("done!")
	if len(errors) > 0 {
		glog.Errorf("%d errors:\n", len(errors))
		for _, err := range errors {
			glog.Error(err)
		}
	}
	if err = dict.Save(); err != nil {
		glog.Fatal(err)
	}
}
