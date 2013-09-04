package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"github.com/hnakamur/commango/os/executil"
)

func main() {
	c := exec.Command("./a.sh")

	stdoutLogger := executil.NewLogger(executil.Info)
	var outBuf bytes.Buffer
	c.Stdout = io.MultiWriter(&outBuf, stdoutLogger)

	stderrLogger := executil.NewLogger(executil.Err)
	var errBuf bytes.Buffer
	c.Stderr = io.MultiWriter(&errBuf, stderrLogger)

	stdoutLogger.Logf("run command\tcommand:%s", executil.CommandLine(c))

	okExitStatuses := []int{0, 1}
	exitStatus, err := executil.Run(c, okExitStatuses)

	// When ok exit status is just zero, you can do:
	// exitStatus, err := executil.Run(c, nil)

	if err != nil {
		stdoutLogger.Logf("failed\tstatus:%d", exitStatus)
	} else {
		stdoutLogger.Logf("done\tstatus:%d", exitStatus)
	}
	fmt.Printf("stdout:%s\n", outBuf.String())
	fmt.Printf("stderr:%s\n", errBuf.String())
	if err != nil {
		panic(err)
	}
}

