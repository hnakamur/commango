package main

import (
	"bytes"
	"os/exec"
	"strings"

	log "github.com/cihub/seelog"
	"github.com/hnakamur/commango/os/executil"
)

func main() {
	configLogger()

	defer log.Flush()
	log.Info("Commango start!")

	c := exec.Command("./a.sh")
	var outBuf bytes.Buffer
	c.Stdout = &outBuf
	var errBuf bytes.Buffer
	c.Stderr = &errBuf

	log.Infof("run command\tcommand:%s", executil.CommandLine(c))
	err := c.Run()
	exitStatus := executil.GetExitStatus(err)
	failed := exitStatus != 0 // && exitStatus != 1
	if failed {
		InfofLines("%s\tout:stdout", outBuf.String())
		WarnfLines("%s\tout:stderr", errBuf.String())
		log.Errorf("failed\trc:%d", exitStatus)
	} else {
		InfofLines("%s\tout:stdout", outBuf.String())
		InfofLines("%s\tout:stderr", errBuf.String())
		log.Infof("done\trc:%d", exitStatus)
	}
	log.Info("Commango finished!")
}

func WarnfLines(format, text string) {
	lines := strings.Split(text, "\n")
	for _, line := range(lines) {
		log.Warnf(format, line)
	}
}

func InfofLines(format, text string) {
	lines := strings.Split(text, "\n")
	for _, line := range(lines) {
		log.Infof(format, line)
	}
}

func configLogger() {
	config := `
<seelog type="sync">
	<outputs>
		<filter levels="trace,debug,info">
			<console formatid="ltsv"/>
		</filter>
		<filter levels="warn,error,critical">
			<console formatid="ltsv_error"/>
		</filter>
		<file formatid="ltsv" path="result.log"/>
	</outputs>
	<formats>
		<format id="ltsv" format="time:%Date(2006-01-02T15:04:05.000Z07:00)%tlev:%l%tmsg:%Msg%n"/>
		<format id="ltsv_error"
			format="%EscM(31)time:%Date(2006-01-02T15:04:05.000Z07:00)%tlev:%l%tmsg:%Msg%EscM(0)%n"/>
	</formats>
</seelog>`

	logger, err := log.LoggerFromConfigAsBytes([]byte(config))
	if err != nil {
		panic(err)
	}
	log.ReplaceLogger(logger)
}
