package main

import (
    "os"

	log "github.com/cihub/seelog"
	"github.com/hnakamur/commango/modules"
//	"github.com/hnakamur/commango/modules/command"
//	"github.com/hnakamur/commango/modules/file"
//	"github.com/hnakamur/commango/modules/directory"
//	"github.com/hnakamur/commango/modules/redhat/service"
//	"github.com/hnakamur/commango/modules/redhat/yum"
	"github.com/hnakamur/commango/modules/shell"
	"github.com/hnakamur/commango/task"
)

func configLogger() {
	config := `
<seelog type="sync">
	<outputs>
		<filter levels="trace">
			<console formatid="skipped"/>
		</filter>
		<filter levels="debug">
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
		<format id="error" format="%EscM(31)%Msg%EscM(0)%n"/>
		<format id="plain" format="%Msg%n"/>
		<format id="unchanged" format="%EscM(32)%Msg%EscM(0)%n"/>
		<format id="skipped" format="%EscM(34)%Msg%EscM(0)%n"/>
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

    queue := task.NewTaskQueue()
    queue.Add(
        &shell.Shell{
            Command: "echo hostname=`hostname`",
        },
        &shell.Shell{
            Chdir: "/tmp",
            Command: "pwd",
        },
    )
    err := queue.RunLoop()
    if err != nil {
        os.Exit(1)
    }

//    command.Command("echo", "a", "b c")
//	directory.Exists("/tmp/foo/bar")
//	//directory.EnsureExists("/tmp/foo/bar", 0755)
//	file.Chown("/tmp/foo", "vagrant:vagrant", true)
//	file.Chmod("/tmp/foo", 0755, true)
//	//directory.EnsureRemoved("/tmp/foo")
//	yum.EnsureInstalled("ntp")
//	service.EnsureStarted("ntpd")
//	service.EnsureAutoStartEnabled("ntpd")
}
