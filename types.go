package d

import (
	"github.com/liuzl/dict"
	"github.com/liuzl/store"
)

type DictValue struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type Values []*DictValue

type Dictionary struct {
	dir   string
	kv    *store.LevelStore
	cedar *dict.Cedar
}
