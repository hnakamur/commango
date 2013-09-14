package directory

import (
	"errors"
	"fmt"
	"os"

	"github.com/hnakamur/commango/modules"
)

func Exists(path string) (bool, error) {
	extra := make(map[string]interface{})
	extra["directory"] = path
	result := modules.Result{
		Extra: extra,
	}

	result.RecordStartTime()
	fi, err := os.Lstat(path)
	result.RecordEndTime()

	var exists bool
	if err != nil {
		if isNoSuchFileOrDirectory(err) {
			err = nil
			exists = false
		}
	} else {
		if fi.IsDir() {
			exists = true
		} else {
			err = errors.New("something not directory exists or permission denied")
		}
	}

	extra["exists"] = exists
	if err != nil {
		result.Err = err
		result.Failed = true
	}
	result.Log()
	modules.ExitOnError(err)
	return exists, nil
}

func isNoSuchFileOrDirectory(err error) bool {
	if er2, ok := err.(*os.PathError); ok {
		return er2.Err.Error() == "no such file or directory"
	} else {
		return false
	}
}

func EnsureExists(path string, perm os.FileMode) (err error) {
	extra := make(map[string]interface{})
	extra["directory"] = path
	extra["op"] = "create"
	extra["mode"] = fmt.Sprintf("%o", perm)
	result := modules.Result{
		Extra: extra,
	}

	result.RecordStartTime()
	fi, err := os.Lstat(path)
	if err != nil {
		if isNoSuchFileOrDirectory(err) {
			err = os.MkdirAll(path, perm)
			if err == nil {
				result.Changed = true
				extra["msg"] = "created"
			}
		}
	} else {
		if fi.IsDir() {
			origMode := fi.Mode() & os.ModePerm
			if origMode != perm {
				extra["orig_mode"] = fmt.Sprintf("%o", origMode)
				err = os.Chmod(path, perm)
				if err == nil {
					result.Changed = true
					extra["msg"] = "directory existed, changed only mode"
				}
			}
		} else {
			err = errors.New("something not directory exists or permission denied")
		}
	}
	result.RecordEndTime()

	if err != nil {
		result.Err = err
		result.Failed = true
	}
	result.Log()
	modules.ExitOnError(err)
	return err
}

func EnsureRemoved(path string) (err error) {
	extra := make(map[string]interface{})
	extra["directory"] = path
	extra["op"] = "remove"
	result := modules.Result{
		Extra: extra,
	}

	result.RecordStartTime()
	fi, err := os.Lstat(path)
	if err != nil {
		if isNoSuchFileOrDirectory(err) {
			err = nil
		}
	} else {
		if fi.IsDir() {
			err = os.RemoveAll(path)
			if err == nil {
				extra["msg"] = "removed"
				result.Changed = true
			}
		} else {
			err = errors.New("something not directory exists or permission denied")
		}
	}
	result.RecordEndTime()

	if err != nil {
		result.Err = err
		result.Failed = true
	}
	result.Log()
	modules.ExitOnError(err)
	return err
}
