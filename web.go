package d

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/liuzl/goutil/rest"
)

func (d *Dictionary) RegisterWeb() {
	http.Handle("/api/get", rest.WithLog(d.GetHandler))
	http.Handle("/api/match", rest.WithLog(d.PrefixMatchHandler))
	http.Handle("/api/update", rest.WithLog(d.UpdateHandler))
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
