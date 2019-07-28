package d

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/liuzl/store"
	"github.com/mitchellh/hashstructure"
)

type Pos struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type Matches struct {
	Value map[string]interface{} `json:"value"`
	Hits  []*Pos                 `json:"hits"`
}

type wordValue struct {
	Typ string
	Val interface{}
}

type pair struct {
	word  string
	value *wordValue
}

func (d *Dictionary) Get(key string) (map[string]interface{}, error) {
	b, err := d.kv.Get(key)
	if err != nil {
		return nil, err
	}
	var v map[string]interface{}
	if err = store.BytesToObject(b, &v); err != nil {
		return nil, err
	}
	return v, nil
}

func (d *Dictionary) PrefixMatch(text string) (
	map[string]map[string]interface{}, error) {

	set := make(map[uint64]*pair)
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
		// type of v is map[string]interface{}
		for typ, val := range v {
			wv := &wordValue{typ, val}
			key, _ := hashstructure.Hash(wv, nil)
			if set[key] == nil || len(word) > len(set[key].word) {
				set[key] = &pair{word, wv}
			}
		}
	}
	ret := make(map[string]map[string]interface{})
	for _, item := range set {
		if ret[item.word] == nil {
			ret[item.word] = map[string]interface{}{item.value.Typ: item.value.Val}
		} else {
			ret[item.word][item.value.Typ] = item.value.Val
		}
	}
	return ret, nil
}

func (d *Dictionary) MultiMatch(text string) (map[string]*Matches, error) {
	r := []rune(text)
	ret := make(map[string]*Matches)
	for i := 0; i < len(r); i++ {
		start := len(string(r[:i]))
		hit, err := d.PrefixMatch(string(r[i:]))
		if err != nil {
			return nil, err
		}
		for k, v := range hit {
			if ret[k] == nil {
				ret[k] = &Matches{v, nil}
			}
			ret[k].Hits = append(ret[k].Hits, &Pos{start, start + len(k)})
		}
	}
	return ret, nil
}

func (d *Dictionary) MultiMaxMatch(text string) (map[string]*Matches, error) {
	r := []rune(text)
	ret := make(map[string]*Matches)
	for i := 0; i < len(r); {
		start := len(string(r[:i]))
		hit, err := d.PrefixMatch(string(r[i:]))
		if err != nil {
			return nil, err
		}
		if len(hit) == 0 {
			i++
			continue
		}
		k := ""
		for key, _ := range hit {
			if len(key) > len(k) {
				k = key
			}
		}
		if ret[k] == nil {
			ret[k] = &Matches{hit[k], nil}
		}
		ret[k].Hits = append(ret[k].Hits, &Pos{start, start + len(k)})
		i += len([]rune(k))
	}
	return ret, nil
}

func (d *Dictionary) Update(k string, values map[string]interface{}) error {
	if k = strings.TrimSpace(k); k == "" {
		return fmt.Errorf("empty key")
	}
	if values == nil || len(values) == 0 {
		return fmt.Errorf("nil values")
	}

	b, err := d.kv.Get(k)
	old := make(map[string]interface{})
	if err == nil {
		if err = store.BytesToObject(b, &old); err != nil {
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

func (d *Dictionary) Replace(k string, values map[string]interface{}) error {
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
	if err = d.kv.Put(k, b); err != nil {
		d.cedar.SafeDelete([]byte(k))
		return err
	}
	atomic.AddInt64(&d.changed, 1)
	atomic.StoreInt64(&d.updated, time.Now().Unix())
	return nil
}

func (d *Dictionary) Delete(k string) error {
	if k = strings.TrimSpace(k); k == "" {
		return fmt.Errorf("empty key")
	}
	if err := d.cedar.SafeDelete([]byte(k)); err != nil {
		return err
	}
	atomic.AddInt64(&d.changed, 1)
	atomic.StoreInt64(&d.updated, time.Now().Unix())
	return d.kv.Delete(k)
}
