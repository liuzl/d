package d

import (
	"github.com/liuzl/dict"
	"github.com/liuzl/store"
)

type Dictionary struct {
	dir   string
	kv    *store.LevelStore
	cedar *dict.Cedar

	changed int64
	updated int64
}
