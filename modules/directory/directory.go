package directory

import (
	"errors"
	"fmt"
	"os"

	"github.com/hnakamur/commango/modules/file"
	"github.com/hnakamur/commango/task"
)

type State string

const (
	Present = State("present")
	Absent  = State("absent")
)

type Directory struct {
	State State
	Path  string
	Owner string
	Group string
	Mode  os.FileMode
}

func (d *Directory) Run() (result *task.Result, err error) {
	result, err = task.DoRun(func(result *task.Result) (err error) {
		result.Module = "directory"
		result.Extra["state"] = string(d.State)
		result.Extra["path"] = d.Path
		if d.State == Absent {
			result.Op = "remove"
			result, err = d.ensureAbsent(result)
		} else {
			result.Op = "create"
			result.Extra["mode"] = fmt.Sprintf("%o", d.Mode)
			result, err = d.ensurePresent(result)
		}
		return
	})

	if d.Owner != "" || d.Group != "" {
		chown := &file.Chown{
			Path:      d.Path,
			Owner:     d.Owner,
			Group:     d.Group,
			Recursive: false,
		}
		result, err = chown.Run()
	}
	return
}

func (d *Directory) ensurePresent(result *task.Result) (*task.Result, error) {
	fi, err := os.Lstat(d.Path)
	if err != nil {
		if isNoSuchFileOrDirectory(err) {
			err = os.MkdirAll(d.Path, d.Mode)
			if err == nil {
				result.Changed = true
			}
		}
	} else {
		if fi.IsDir() {
			oldMode := fi.Mode() & os.ModePerm
			if oldMode == d.Mode {
				result.Skipped = true
			} else {
				result.Extra["old_mode"] = fmt.Sprintf("%o", oldMode)
				err = os.Chmod(d.Path, d.Mode)
				if err == nil {
					result.Changed = true
				}
			}
		} else {
			err = errors.New("something not directory exists or permission denied")
		}
	}
	return result, err
}

func (d *Directory) ensureAbsent(result *task.Result) (*task.Result, error) {
	fi, err := os.Lstat(d.Path)
	if err != nil {
		if isNoSuchFileOrDirectory(err) {
			result.Skipped = true
			err = nil
		}
	} else {
		if fi.IsDir() {
			err = os.RemoveAll(d.Path)
			if err == nil {
				result.Changed = true
			}
		} else {
			err = errors.New("something not directory exists or permission denied")
		}
	}
	return result, err
}

func isNoSuchFileOrDirectory(err error) bool {
	if er2, ok := err.(*os.PathError); ok {
		return er2.Err.Error() == "no such file or directory"
	} else {
		return false
	}
}
