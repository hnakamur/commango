package template

import (
	"fmt"
	"os"

	"github.com/hnakamur/commango/modules"
	"github.com/hnakamur/commango/templateutil"
)

func EnsureExists(path string, content string, data interface{},
	perm os.FileMode) (result modules.Result, err error) {
	result.RecordStartTime()

	extra := map[string]interface{}{
		"op":      "template",
		"path":    path,
		"content": content,
		"data":    data,
		"mode":    fmt.Sprintf("%o", perm),
	}

	defer func() {
		result.Extra = extra

		result.RecordEndTime()

		if err != nil {
			result.Err = err
			result.Failed = true
		}
		result.Log()
		modules.ExitOnError(err)
	}()

	tmpl, err := templateutil.NewWithString(path, content)
	if err != nil {
		return
	}

	changed, err := templateutil.WriteIfChanged(tmpl, data, path, perm)
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
