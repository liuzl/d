package main

import (
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/liuzl/d"
)

var (
	dir  = flag.String("dir", "data", "data dir")
	addr = flag.String("addr", ":8080", "band address")
)

func main() {
	flag.Parse()
	defer glog.Flush()

	dict, err := d.Load(*dir)
	if err != nil {
		glog.Fatal(err)
	}

	dict.RegisterWeb()

	glog.Info("dserver listen on", *addr)
	glog.Error(http.ListenAndServe(":8080", nil))
}
