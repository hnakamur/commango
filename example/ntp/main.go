package main

import (
	"os"

	log "github.com/cihub/seelog"
	"github.com/hnakamur/commango/modules"
	//	"github.com/hnakamur/commango/modules/command"
	//	"github.com/hnakamur/commango/modules/file"
	"github.com/hnakamur/commango/modules/directory"
	"github.com/hnakamur/commango/modules/template"
	//	"github.com/hnakamur/commango/modules/redhat/service"
	"github.com/hnakamur/commango/modules/redhat/group"
	"github.com/hnakamur/commango/modules/redhat/user"
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

const NTP_CONF_TEMPLATE = `# For more information about this file, see the man pages
# ntp.conf(5), ntp_acc(5), ntp_auth(5), ntp_clock(5), ntp_misc(5), ntp_mon(5).

driftfile /var/lib/ntp/drift

# Permit time synchronization with our time source, but do not
# permit the source to query or modify the service on this system.
restrict default kod nomodify notrap nopeer noquery
restrict -6 default kod nomodify notrap nopeer noquery

# Permit all access over the loopback interface.  This could
# be tightened as well, but to do so would effect some of
# the administrative functions.
restrict 127.0.0.1 
restrict -6 ::1

# Hosts on local network are less restricted.
#restrict 192.168.1.0 mask 255.255.255.0 nomodify notrap

# Use public servers from the pool.ntp.org project.
# Please consider joining the pool (http://www.pool.ntp.org/join.html).
{{range .ntp_servers}}{{/*
*/}}server {{.}}
{{end}}
#broadcast 192.168.1.255 autokey    # broadcast server
#broadcastclient            # broadcast client
#broadcast 224.0.1.1 autokey        # multicast server
#multicastclient 224.0.1.1      # multicast client
#manycastserver 239.255.254.254     # manycast server
#manycastclient 239.255.254.254 autokey # manycast client

# Undisciplined Local Clock. This is a fake driver intended for backup
# and when no outside source of synchronized time is available. 
#server 127.127.1.0 # local clock
#fudge  127.127.1.0 stratum 10  

# Enable public key cryptography.
#crypto

includefile /etc/ntp/crypto/pw

# Key file containing the keys and key identifiers used when operating
# with symmetric key cryptography. 
keys /etc/ntp/keys

# Specify the key identifiers which are trusted.
#trustedkey 4 8 42

# Specify the key identifier to use with the ntpdc utility.
#requestkey 8

# Specify the key identifier to use with the ntpq utility.
#controlkey 8

# Enable writing of statistics records.
#statistics clockstats cryptostats loopstats peerstats
`

func main() {
	configLogger()
	modules.EnableExitOnError()

	queue := task.NewTaskQueue()
	queue.Add(
		&shell.Shell{
			Command: "echo hostname=`hostname`",
		},
		&shell.Shell{
			Chdir:   "/tmp",
			Command: "pwd",
		},
		&directory.Directory{
			State: directory.Present,
			Path:  "/tmp/foo/bar",
			Mode:  0755,
		},
		&template.Template{
			Path:    "/tmp/foo/bar/baz.conf",
			Content: NTP_CONF_TEMPLATE,
			Data: map[string]interface{}{
				"ntp_servers": []string{
					"ntp.nict.jp",
					"ntp.jst.mfeed.ad.jp",
					"ntp.ring.gr.jp",
				},
			},
			Mode: 0644,
		},
		&user.User{
			State: user.Present,
			Name:  "foo",
			Uid:   user.AUTO_UID,
		},
		&group.Group{
			State: group.Present,
			Name:  "bar",
			Gid:   group.AUTO_GID,
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
