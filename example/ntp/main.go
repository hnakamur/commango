package main

import (
	log "github.com/cihub/seelog"
	"github.com/hnakamur/commango/modules"
	"github.com/hnakamur/commango/modules/file"
	"github.com/hnakamur/commango/modules/directory"
	"github.com/hnakamur/commango/modules/redhat/service"
	"github.com/hnakamur/commango/modules/redhat/yum"
)

func configLogger() {
	config := `
<seelog type="sync">
	<outputs>
		<filter levels="trace,debug">
			<console formatid="unchanged"/>
		</filter>
		<filter levels="info">
			<console formatid="plain"/>
		</filter>
		<filter levels="warn,error,critical">
			<console formatid="error"/>
		</filter>
	</outputs>
	<formats>
		<format id="plain" format="%Msg%n"/>
		<format id="unchanged" format="%EscM(32)%Msg%EscM(0)%n"/>
		<format id="error" format="%EscM(31)%Msg%EscM(0)%n"/>
	</formats>
</seelog>`

	logger, err := log.LoggerFromConfigAsBytes([]byte(config))
	if err != nil {
		panic(err)
	}
	log.ReplaceLogger(logger)
}

func main() {
	configLogger()
	modules.EnableExitOnError()

	directory.Exists("/tmp/foo/bar")
	directory.EnsureExists("/tmp/foo/bar", 0755)
	file.Chown("/tmp/foo", "vagrant:vagrant", true)
	file.Chmod("/tmp/foo", 0755, true)
	//directory.EnsureRemoved("/tmp/foo")
	yum.EnsureInstalled("ntp")
	service.EnsureStarted("ntpd")
	service.EnsureAutoStartEnabled("ntpd")
}
