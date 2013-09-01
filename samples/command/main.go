package main

import (
	"fmt"
	"os/exec"
	"github.com/hnakamur/commango/os/executil"
)

//func back() {
//	cmd := exec.Command("./a.sh")
//	var out bytes.Buffer
//	cmd.Stdout = &out
//	err := cmd.Run()
//	if err != nil {
//		if e2, ok := err.(*exec.ExitError); ok {
//			if s, ok := e2.Sys().(syscall.WaitStatus); ok {
//				Logf(Err, "exit_code:%d", int(s.ExitStatus()))
//			}
//		} else {
//			panic(err)
//		}
//	}
//	Logf(Info, "%s", out.String())
//	Logf(Err, "%s", out.String())
//}
//
//func main2() {
//	cmd := exec.Command("./a.sh", ">", "/tmp/hoge") // Not redirected.
//	//cmd := exec.Command("sh", "-c", "./a.sh > /tmp/hoge")
//	cmd.Stdout = NewLogger(Info)
//	cmd.Stderr = NewLogger(Err)
//	err := cmd.Run()
//	if err != nil {
//		if e2, ok := err.(*exec.ExitError); ok {
//			if s, ok := e2.Sys().(syscall.WaitStatus); ok {
//				cmd.Stderr.Write([]byte(fmt.Sprintf("failed\texit_code:%d",
//					int(s.ExitStatus()))))
//			} else {
//				panic(errors.New("Unimplemented for system where exec.ExitError.Sys() is not syscall.WaitStatus."))
//			}
//		} else {
//			panic(err)
//		}
//	}
//}

func main() {
	runner := executil.CommandRunner{
		//Command: exec.Command("uname", "-a"),
		Command: exec.Command("./a.sh"),
		OkExitStatuses: []int{0, 1},
		CapturesStdout: true,
		CapturesStderr: true,
	}
	err := runner.Run()
	fmt.Printf("stdout:%s\n", runner.StdoutOutput())
	fmt.Printf("stderr:%s\n", runner.StderrOutput())
	if err != nil {
		panic(err)
	}
}

