package d

import (
	"fmt"
	"strings"

	"github.com/liuzl/store"
)

func (d *Dictionary) Get(key string) (Values, error) {
	b, err := d.kv.Get(key)
	if err != nil {
		return nil, err
	}
	var v Values
	if err = store.BytesToObject(b, &v); err != nil {
		return nil, err
	}
	return v, nil
}

func (d *Dictionary) PrefixMatch(text string) (map[string]Values, error) {
	ret := make(map[string]Values)
	for _, id := range d.cedar.PrefixMatch([]byte(text), 0) {
		key, err := d.cedar.Key(id)
		if err != nil {
			return nil, err
		}
		word := string(key)
		v, err := d.Get(word)
		if err != nil {
			return nil, err
		}
		ret[word] = v
	}
	return ret, nil
}

func (d *Dictionary) Update(k string, values Values) error {
	if k = strings.TrimSpace(k); k == "" {
		return fmt.Errorf("empty key")
	}
	if values == nil || len(values) == 0 {
		return fmt.Errorf("nil values")
	}

	b, err := d.kv.Get(k)
	old := make(Values)
	if err == nil {
		if err = store.BytesToObject(b, old); err != nil {
			return err
		}
	} else if err.Error() != "leveldb: not found" {
		return err
	}
	for k, v := range values {
		old[k] = v
	}
	return d.Replace(k, old)
}

func (d *Dictionary) Replace(k string, values Values) error {
	if k = strings.TrimSpace(k); k == "" {
		return fmt.Errorf("empty key")
	}
	if values == nil || len(values) == 0 {
		return fmt.Errorf("nil values")
	}

	var b []byte
	var err error
	if err = d.cedar.SafeInsert([]byte(k), len(values)); err != nil {
		return err
	}
	if b, err = store.ObjectToBytes(values); err != nil {
		d.cedar.SafeDelete([]byte(k))
		return err
	}
	return d.kv.Put(k, b)
}

func (d *Dictionary) Delete(k string) error {
	if k = strings.TrimSpace(k); k == "" {
		return fmt.Errorf("empty key")
	}
	if err := d.cedar.SafeDelete([]byte(k)); err != nil {
		return err
	}
	return d.kv.Delete(k)
}
