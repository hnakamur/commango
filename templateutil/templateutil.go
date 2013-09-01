package templateutil

import (
	"bytes"
	"text/template"
)

func RenderToBytes(template *template.Template, data interface{}) (output []byte, err error) {
	var buf bytes.Buffer
	err = template.Execute(&buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
