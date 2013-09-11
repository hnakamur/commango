package main

import (
	"os/exec"
	"strings"

	log "github.com/cihub/seelog"
	"github.com/hnakamur/commango/os/executil"
)

func main() {
	configLogger()

	defer log.Flush()
	log.Info("Commango start!")
	//RunCommand("./a.sh")
	RunCommand("touch", "foo bar")
	log.Info("Commango finished!")
}

func RunCommand(name string, arg ...string) bool {
	cmd := exec.Command(name, arg...)
	s, err := executil.FormatCommand(cmd)
	if err != nil {
		log.Errorf("failed\terr:%s", err)
		return true
	}
	log.Infof("run command\tcommand:%s", s)
	r, err := executil.Run(cmd)
	failed := r.Rc != 0 && r.Rc != 1
	if err != nil && !executil.IsExitError(err) {
		InfofLines("%s\tout:stdout", r.Out.String())
		WarnfLines("%s\tout:stderr", r.Err.String())
		log.Errorf("failed\terr:%s", err)
	} else if failed {
		InfofLines("%s\tout:stdout", r.Out.String())
		WarnfLines("%s\tout:stderr", r.Err.String())
		log.Errorf("failed\trc:%d\terr:%s", r.Rc, err)
	} else {
		InfofLines("%s\tout:stdout", r.Out.String())
		InfofLines("%s\tout:stderr", r.Err.String())
		log.Infof("done\trc:%d", r.Rc)
	}
	return failed
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
