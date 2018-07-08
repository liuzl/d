package d

import (
	"github.com/liuzl/dict"
	"github.com/liuzl/store"
)

type Values map[string]interface{}

type Dictionary struct {
	dir   string
	kv    *store.LevelStore
	cedar *dict.Cedar
}
