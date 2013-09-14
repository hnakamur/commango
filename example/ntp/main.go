package main

import (
	log "github.com/cihub/seelog"
	"github.com/hnakamur/commango/modules/redhat/service"
	"github.com/hnakamur/commango/modules/redhat/yum"
)

func configLogger() {
	config := `
<seelog type="sync">
	<outputs>
		<filter levels="trace,debug,info">
			<console formatid="plain"/>
		</filter>
		<filter levels="warn,error,critical">
			<console formatid="error"/>
		</filter>
	</outputs>
	<formats>
		<format id="plain" format="%EscM(32)%Msg%EscM(0)%n"/>
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

	installed, err := yum.Installed("ntp")
	if err != nil {
		panic(err)
	}

	if !installed {
		_, err = yum.Install("ntp")
		if err != nil {
			panic(err)
		}
	}

	status, err := service.Status("ntpd")
	if err != nil {
		panic(err)
	}
	if status == service.STOPPED {
		_, err = service.Start("ntpd")
		if err != nil {
			panic(err)
		}
	}

	enabled, err := service.AutoStartEnabled("ntpd")
	if err != nil {
		panic(err)
	}
	if !enabled {
		_, err = service.EnableAutoStart("ntpd")
		if err != nil {
			panic(err)
		}
	}
}
