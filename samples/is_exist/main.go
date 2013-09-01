package main

import (
	"fmt"
	"github.com/hnakamur/commango/os/osutil"
)

func main() {
	fmt.Printf("main.go %v\n", osutil.IsExist("main.go"))
	fmt.Printf("sub.go %v\n", osutil.IsExist("sub.go"))
}
