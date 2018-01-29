package htmlformfill

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

const nameField = "name"
const valField = "value"

// Fill accepts an io.Reader for the HTML and map[string]string of the fields to be set
// Parse html document, filling all fields as required
// Returns filled html document as []byte
func Fill(r io.Reader, f map[string]string) (*bytes.Reader, error) {
	var out []byte
	var e error
	z := html.NewTokenizer(r)
	for {
		t := z.Next()
		if t == html.ErrorToken {
			o := bytes.NewReader(out)
			return o, e
		}
		n, _ := z.TagName()
		o := z.Raw()
		switch string(n) {
		case "input":
			o = input(z, f)
		case "textarea":
			o = textarea(z, f)
		case "select":
			o = selector(z, f)
		}
		out = append(out, o...)
	}
}

// input handles text, number, email, etc. inputs
// Returns filled input row if found in f, original row if not found
func input(z *html.Tokenizer, f map[string]string) []byte {
	var out []byte
	for {
		key, val, _ := z.TagAttr()
		if string(key) == "type" && string(val) == "radio" {
			out = radio(z, f)
			break
		} else if string(key) == "type" && string(val) == "checkbox" {
			out = checkbox(z, f)
			break
		} else if kv, ok := f[string(val)]; string(key) == nameField && ok {
			r := string(z.Raw())
			r = strings.Replace(r, ">", fmt.Sprintf(" value=\"%s\">", kv), -1)
			out = []byte(r)
			break
		} else {
			out = z.Raw()
		}
	}
	return out
}

// radio handles radio inputs
// Returns filled input row if found in f, original row if not found
func radio(z *html.Tokenizer, f map[string]string) []byte {
	var out []byte
	var n string
	for {
		key, val, m := z.TagAttr()
		if string(key) == nameField {
			n = string(val)
		}
		if string(key) == valField && f[n] == string(val) {
			r := string(z.Raw())
			r = strings.Replace(r, ">", " checked>", -1)
			out = []byte(r)
			break
		}
		if !m {
			break
		}
	}
	if len(out) == 0 {
		out = z.Raw()
	}
	return out
}

// checkbox handles checkbox inputs
// Returns filled input row if found in f, original row if not found
func checkbox(z *html.Tokenizer, f map[string]string) []byte {
	var out []byte
	var n string
	for {
		key, val, m := z.TagAttr()
		if string(key) == nameField {
			n = string(val)
		}
		sel := strings.Split(f[n], ",")
		if string(key) == valField {
			for _, v := range sel {
				if string(val) == v {
					r := string(z.Raw())
					r = strings.Replace(r, ">", " checked>", -1)
					out = []byte(r)
				}
			}
		}
		if !m {
			break
		}
	}
	if len(out) == 0 {
		out = z.Raw()
	}
	return out
}

// textarea handles textarea inputs
// Returns filled textarea if found in f, original row if not found
func textarea(z *html.Tokenizer, f map[string]string) []byte {
	var out []byte
	for {
		key, val, m := z.TagAttr()
		if kv, ok := f[string(val)]; string(key) == nameField && ok {
			r := string(z.Raw())
			r = strings.Replace(r, ">", ">"+kv, 1)
			out = []byte(r)
		} else {
			out = z.Raw()
		}
		if !m {
			break
		}
	}
	return out
}

// selector handles select inputs
// Returns filled textarea if found in f, original select block if not found
func selector(z *html.Tokenizer, f map[string]string) []byte {
	var out []byte
	var n string
	for {
		key, val, m := z.TagAttr()
		if _, ok := f[string(val)]; string(key) == nameField && ok {
			n = string(val)
			out = append(out, z.Raw()...)
		}
		if !m {
			break
		}
	}
	for {
		z.Next()
		for {
			key, val, m := z.TagAttr()
			if kv, ok := f[string(n)]; string(key) == valField && ok {
				if kv == string(val) {
					r := string(z.Raw())
					r = strings.Replace(r, ">", " selected>", 1)
					out = append(out, []byte(r)...)
				} else {
					out = append(out, z.Raw()...)
				}
			} else {
				out = append(out, z.Raw()...)
			}
			if !m {
				break
			}
		}
		n, _ := z.TagName()
		if string(n) == "select" {
			break
		}
	}
	return out
}
