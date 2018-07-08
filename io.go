package d

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb"
	"github.com/liuzl/dict"
	"github.com/liuzl/goutil"
	"github.com/liuzl/store"
)

type Record struct {
	K string `json:"k"`
	V Values `json:"v"`
}

func Load(dir string) (*Dictionary, error) {
	kvDir := filepath.Join(dir, "kv")
	cedarDir := filepath.Join(dir, "cedar")
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
	d := &Dictionary{dir, kv, cedar}
	return d, nil
}

func (d *Dictionary) Save() error {
	cedarDir := filepath.Join(d.dir, "cedar")
	return d.cedar.SaveToFile(cedarDir, "gob")
}

func (d *Dictionary) Dump(path string) error {
	wf, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer wf.Close()
	err = d.kv.ForEach(nil, func(key, value []byte) (bool, error) {
		var v Values
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
