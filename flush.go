package d

import (
	"time"
)

func (d *Dictionary) flush() {
	for {
		time.Sleep(5 * time.Second)
		if d.changed > 1000 || time.Now().Unix()-d.updated > 60*5 {
			d.Save()
		}
	}
}
