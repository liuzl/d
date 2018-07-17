package d

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/liuzl/goutil/rest"
)

func (d *Dictionary) RegisterWeb() {
	http.Handle(fmt.Sprintf("/%s/get", d.Name), rest.WithLog(d.GetHandler))
	http.Handle(fmt.Sprintf("/%s/match", d.Name),
		rest.WithLog(d.PrefixMatchHandler))
	http.Handle(fmt.Sprintf("/%s/multimatch", d.Name),
		rest.WithLog(d.MultiMatchHandler))
	http.Handle(fmt.Sprintf("/%s/update", d.Name),
		rest.WithLog(d.UpdateHandler))
}

func (d *Dictionary) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	word := strings.TrimSpace(r.FormValue("word"))
	values, err := d.Get(word)
	if err != nil {
		rest.MustEncode(w, &rest.RestMessage{"ERROR", err.Error()})
		return
	}
	rest.MustEncode(w, &rest.RestMessage{"OK", values})
}

func (d *Dictionary) PrefixMatchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	text := strings.TrimSpace(r.FormValue("text"))
	ret, err := d.PrefixMatch(text)
	if err != nil {
		rest.MustEncode(w, &rest.RestMessage{"ERROR", err.Error()})
		return
	}
	rest.MustEncode(w, &rest.RestMessage{"OK", ret})
}

func (d *Dictionary) MultiMatchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	text := strings.TrimSpace(r.FormValue("text"))
	ret, err := d.MultiMatch(text)
	if err != nil {
		rest.MustEncode(w, &rest.RestMessage{"ERROR", err.Error()})
		return
	}
	rest.MustEncode(w, &rest.RestMessage{"OK", ret})
}

func (d *Dictionary) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	replace := strings.TrimSpace(r.FormValue("replace"))
	flush := strings.TrimSpace(r.FormValue("flush"))
	data := strings.TrimSpace(r.FormValue("json"))
	var rec Record
	if err := json.Unmarshal([]byte(data), &rec); err != nil {
		rest.MustEncode(w, &rest.RestMessage{"ERROR", err.Error()})
		return
	}
	f := d.Update
	if replace != "" {
		f = d.Replace
	}
	if err := f(rec.K, rec.V); err != nil {
		rest.MustEncode(w, &rest.RestMessage{"ERROR", err.Error()})
		return
	}
	if flush != "" {
		if err := d.Save(); err != nil {
			rest.MustEncode(w, &rest.RestMessage{"ERROR", err.Error()})
			return
		}
	}
	rest.MustEncode(w, &rest.RestMessage{"OK", "done"})
}
