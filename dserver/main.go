package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/liuzl/d"
)

var (
	dir = flag.String("dir", "data", "data dir")
)

func main() {
	flag.Parse()
	dict, err := d.Load(*dir)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dict)
	dict.RegisterWeb()
	http.ListenAndServe(":8080", nil)
}
