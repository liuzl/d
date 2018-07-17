package d

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/cheggaaa/pb"
	"github.com/liuzl/dict"
	"github.com/liuzl/goutil"
	"github.com/liuzl/store"
)

var (
	dir = flag.String("dict_dir", "dicts", "dir for dicts")
)

type Record struct {
	K string                 `json:"k"`
	V map[string]interface{} `json:"v"`
}

func Load(name string) (*Dictionary, error) {
	kvDir := filepath.Join(*dir, name, "kv")
	cedarDir := filepath.Join(*dir, name, "cedar")
	kv, err := store.NewLevelStore(kvDir)
	if err != nil {
		return nil, err
	}
	cedar := dict.New()
	if _, err = os.Stat(cedarDir); err == nil {
		err = cedar.LoadFromFile(cedarDir, "gob")
		if err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	d := &Dictionary{Name: name, kv: kv, cedar: cedar}
	go d.flush()
	return d, nil
}

func (d *Dictionary) Save() error {
	cedarDir := filepath.Join(*dir, d.Name, "cedar")
	fmt.Println(cedarDir, "saving")
	if err := d.cedar.SaveToFile(cedarDir, "gob"); err != nil {
		return err
	}
	atomic.StoreInt64(&d.changed, 0)
	atomic.StoreInt64(&d.updated, 0)
	return nil
}

func (d *Dictionary) Dump(path string) error {
	wf, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer wf.Close()
	err = d.kv.ForEach(nil, func(key, value []byte) (bool, error) {
		var v map[string]interface{}
		if err := store.BytesToObject(value, &v); err != nil {
			return false, err
		}
		b, err := json.Marshal(&Record{string(key), v})
		if err != nil {
			return false, err
		}
		if _, err = wf.Write(append(b, '\n')); err != nil {
			return false, err
		}
		return true, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func Build(src, dst string) (*Dictionary, error) {
	count, err := goutil.FileLineCount(src)
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		return nil, fmt.Errorf("empty src file: %s", src)
	}
	d, err := Load(dst)
	if err != nil {
		return nil, err
	}
	fd, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	var rec Record
	br := bufio.NewReader(fd)
	bar := pb.StartNew(count)
	for {
		b, c := br.ReadBytes('\n')
		if c == io.EOF {
			break
		}
		if err = json.Unmarshal(b, &rec); err != nil {
			return nil, err
		}
		if err = d.Update(rec.K, rec.V); err != nil {
			return nil, err
		}
		bar.Increment()
	}
	if err = d.Save(); err != nil {
		return nil, err
	}
	return d, nil
}
