package htmlformfill

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestFill(t *testing.T) {
	f := make(map[string]string)
	f["name"] = "foo bar"
	f["phone"] = "5555555"
	f["textblock"] = "this is a longer block of text"
	f["single-radio"] = "test1"
	f["multi-check"] = "test1,test3"
	f["selector"] = "test2"
	r, err := os.Open("./examples/form.html")
	if err != nil {
		t.Error(err)
	}
	out, err := Fill(r, f)
	if err != nil {
		t.Error(err)
	}
	sout, err := ioutil.ReadAll(out)
	if err != nil {
		t.Error(err)
	}
	werr := ioutil.WriteFile("./examples/test-output.html", sout, 0644)
	if werr != nil {
		t.Error(werr)
	}
	for k, v := range f {
		sv := strings.Split(v, ",")
		if len(sv) > 0 {
			for _, i := range sv {
				if !strings.Contains(string(sout), i) {
					t.Error("Error filling field: " + k)
				}
			}
		} else {
			if !strings.Contains(string(sout), v) {
				t.Error("Error filling field: " + k)
			}
		}
	}
}
