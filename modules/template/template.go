package template

import (
	"fmt"
	"os"

	"github.com/hnakamur/commango/task"
	"github.com/hnakamur/commango/templateutil"
)

type Template struct {
	Path    string
	Content string
	Data    interface{}
	Mode    os.FileMode
}

func (t *Template) Run() (result *task.Result, err error) {
	result = task.NewResult("template")
	result.RecordStartTime()
	defer func() {
        result.Err = err
        result.RecordEndTime()
        result.Log()
    }()

	result.Extra["path"] = t.Path
	result.Extra["content"] = t.Content
	result.Extra["data"] = t.Data
	result.Extra["mode"] = fmt.Sprintf("%o", t.Mode)

	tmpl, err := templateutil.NewWithString(t.Path, t.Content)
	if err != nil {
		return
	}

	changed, err := templateutil.WriteIfChanged(tmpl, t.Data, t.Path, t.Mode)
	if err != nil {
		return
	}

	if changed {
		result.Changed = true
	} else {
		result.Skipped = true
	}
	return
}
