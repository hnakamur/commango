package main

import (
	"fmt"
	"github.com/hnakamur/commango/os/osutil"
)

func main() {
	fmt.Printf("main.go %v\n", osutil.Exists("main.go"))
	fmt.Printf("sub.go %v\n", osutil.Exists("sub.go"))
}
