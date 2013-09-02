package templateutil

import (
	"bytes"
	"io/ioutil"
	"os"
	"text/template"
)

func NewWithString(name, content string) (*template.Template, error) {
	tmpl := template.New(name)
	return tmpl.Parse(content)
}

func WriteIfChanged(tmpl *template.Template, data interface{},
		dest_path string, perm os.FileMode) (changed bool, err error) {
	output, err := RenderToBytes(tmpl, data)
	if err != nil {
		return false, err
	}

	orig, err := ioutil.ReadFile(dest_path)
	if err != nil {
		return false, err
	}

	if bytes.Equal(output, orig) {
		return false, nil
	}

	err = ioutil.WriteFile(dest_path, output, perm)
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
