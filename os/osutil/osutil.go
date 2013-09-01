package osutil

import (
	"os"
)

func IsExist(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
