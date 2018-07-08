package d

import (
	"time"
)

func (d *Dictionary) sync() {
	for {
		time.Sleep(5 * time.Second)
		if d.changed > 10 || time.Now().Unix()-d.updated > 60*5 {
			d.Save()
		}
	}
}
