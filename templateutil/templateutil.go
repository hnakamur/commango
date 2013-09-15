package templateutil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

func NewWithString(name, content string) (*template.Template, error) {
	tmpl := template.New(name)
	return tmpl.Parse(content)
}

func WriteIfChanged(tmpl *template.Template, data interface{},
		path string, perm os.FileMode) (changed bool, err error) {
	output, err := RenderToBytes(tmpl, data)
	if err != nil {
		return false, err
	}

	fi, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}

	if fi != nil {
		if !fi.Mode().IsRegular() {
			return false, fmt.Errorf("something not file exists: %s", path)
		}

		orig, err := ioutil.ReadFile(path)
		if err != nil {
			return false, err
		}

		if bytes.Equal(output, orig) {
			return false, nil
		}
	}

	err = ioutil.WriteFile(path, output, perm)
	if err != nil {
		return false, err
	}

	return true, nil
}

func RenderToBytes(tmpl *template.Template, data interface{}) (output []byte, err error) {
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
