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
	f["single-check"] = "test3"
	f["selector"] = "test2"
	r, err := os.Open("./examples/form.html")
	if err != nil {
		t.Error(err)
	}
	out, err := Fill(r, f)
	if err != nil {
		t.Error(err)
	}
	for k, v := range f {
		if !strings.Contains(string(out), v) {
			t.Error("Error filling field: " + k)
		}
	}
	werr := ioutil.WriteFile("./examples/test-output.html", out, 0644)
	if werr != nil {
		t.Error(werr)
	}
}
