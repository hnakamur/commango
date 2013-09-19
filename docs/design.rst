デザインメモ
============

まだフワフワしていて固まっていないです。
ご意見お待ちしています。

基本的な使い方の想定
--------------------

1. 実行したいタスクリストをGoのプログラムとして記述
2. ターゲットマシン用にクロスコンパイル
3. 生成したバイナリ実行モジュールをターゲットマシンにファイル転送
4. ターゲットマシン上の実行モジュールを制御マシンからsshで呼び出して実行

将来的には、転送した後デーモンとして起動してコマンド呼び出しで制御するというのも
考えてもいいかも。

ユーザ視点の希望
----------------

- 実行されたコマンドの内容と出力はログで全て確認したい。

    - chefとかansibleでタスクの内容がシェルスリクプトの一時ファイルを作って対象マシンに転送して実行しているが、そのファイル名をログで見せられても一時ファイルは消されてしまうのであまり嬉しくない。


概念と用語
----------

  タスク (task.Task)
    処理のユーザからの呼び出し単位。インタフェースとして定義。

  モジュール
    各機能をタスクとして実装したもの。冪等性をもたせる。

  実行結果 (task.Result)
    タスクまたはサブタスクの実行結果。コマンドを実行したら必ず実行結果を
    生成してログ出力する。
    現在のモヤモヤポイントは1つのタスクが複数のコマンド実行（サブタスク）で
    構成される場合にTaskのRun()の戻り値は最後のResultだけを返しているところ。

実装メモ
--------

pure GoのソースはGOOS=linux GOARCH=amd64 go buildなどとしてクロスコンパイル
できますが、cgoが必要なソースはできません。インクルードファイルやライブラリ
が必要になるので仕方ない気がします。

そこで、Linuxのユーザやグループを扱うタスクはcgo依存のパッケージを使わずに、
useraddなどの外部コマンドを使う方法で実装しています。

実行例
------

::

  [vagrant@cent64 commango]$ cd example/ntp
  [vagrant@cent64 ntp]$ go build
  [vagrant@cent64 ntp]$ sudo ./ntp 
  {"changed":true,"cmd":"/bin/sh -c \"echo hostname=`hostname`\"","command":"echo hostname=`hostname`","delta":"0:00:00.00264","end":"2013-09-20 01:00:10.02498","failed":false,"module":"shell","rc":0,"shell":"/bin/sh","skipped":false,"start":"2013-09-20 01:00:10.02234","stderr":"","stdout":"hostname=cent64.internal.example.com\n"}
  {"changed":true,"chdir":"/tmp","cmd":"/bin/sh -c \"cd /tmp; pwd\"","command":"pwd","delta":"0:00:00.00164","end":"2013-09-20 01:00:10.02678","failed":false,"module":"shell","rc":0,"shell":"/bin/sh","skipped":false,"start":"2013-09-20 01:00:10.02514","stderr":"","stdout":"/tmp\n"}
  {"changed":true,"delta":"0:00:00.00002","end":"2013-09-20 01:00:10.02690","failed":false,"mode":"755","module":"directory","old_mode":"750","path":"/tmp/foo/bar","skipped":false,"start":"2013-09-20 01:00:10.02688","state":"present"}
  {"changed":false,"cmd":"find /tmp/foo/bar -printf %u:%g\\n","delta":"0:00:00.00111","end":"2013-09-20 01:00:10.02807","failed":false,"module":"chown.get_owners","path":"/tmp/foo/bar","rc":0,"recursive":false,"skipped":false,"start":"2013-09-20 01:00:10.02696","stderr":"","stdout":"root:root\nvagrant:vagrant\n"}
  {"changed":true,"cmd":"chown root:root /tmp/foo/bar","delta":"0:00:00.00072","end":"2013-09-20 01:00:10.02888","failed":false,"group":"root","module":"chown","owner":"root","path":"/tmp/foo/bar","rc":0,"recursive":false,"skipped":false,"start":"2013-09-20 01:00:10.02816","stderr":"","stdout":""}
  {"changed":false,"cmd":"find /tmp/foo -printf %m\\n","delta":"0:00:00.00092","end":"2013-09-20 01:00:10.02989","failed":false,"module":"chmod.get_modes","path":"/tmp/foo","rc":0,"recursive":true,"skipped":false,"start":"2013-09-20 01:00:10.02897","stderr":"","stdout":"750\n755\n750\n"}
  {"changed":true,"cmd":"chmod -R 750 /tmp/foo","delta":"0:00:00.00105","end":"2013-09-20 01:00:10.03102","failed":false,"mode":"750","module":"chmod","old_modes":["750","755"],"path":"/tmp/foo","rc":0,"recursive":true,"skipped":false,"start":"2013-09-20 01:00:10.02997","stderr":"","stdout":""}
  {"changed":false,"content":"# For more information about this file, see the man pages\n# ntp.conf(5), ntp_acc(5), ntp_auth(5), ntp_clock(5), ntp_misc(5), ntp_mon(5).\n\ndriftfile /var/lib/ntp/drift\n\n# Permit time synchronization with our time source, but do not\n# permit the source to query or modify the service on this system.\nrestrict default kod nomodify notrap nopeer noquery\nrestrict -6 default kod nomodify notrap nopeer noquery\n\n# Permit all access over the loopback interface.  This could\n# be tightened as well, but to do so would effect some of\n# the administrative functions.\nrestrict 127.0.0.1 \nrestrict -6 ::1\n\n# Hosts on local network are less restricted.\n#restrict 192.168.1.0 mask 255.255.255.0 nomodify notrap\n\n# Use public servers from the pool.ntp.org project.\n# Please consider joining the pool (http://www.pool.ntp.org/join.html).\n{{range .ntp_servers}}{{/*\n*/}}server {{.}}\n{{end}}\n#broadcast 192.168.1.255 autokey    # broadcast server\n#broadcastclient            # broadcast client\n#broadcast 224.0.1.1 autokey        # multicast server\n#multicastclient 224.0.1.1      # multicast client\n#manycastserver 239.255.254.254     # manycast server\n#manycastclient 239.255.254.254 autokey # manycast client\n\n# Undisciplined Local Clock. This is a fake driver intended for backup\n# and when no outside source of synchronized time is available. \n#server 127.127.1.0 # local clock\n#fudge  127.127.1.0 stratum 10  \n\n# Enable public key cryptography.\n#crypto\n\nincludefile /etc/ntp/crypto/pw\n\n# Key file containing the keys and key identifiers used when operating\n# with symmetric key cryptography. \nkeys /etc/ntp/keys\n\n# Specify the key identifiers which are trusted.\n#trustedkey 4 8 42\n\n# Specify the key identifier to use with the ntpdc utility.\n#requestkey 8\n\n# Specify the key identifier to use with the ntpq utility.\n#controlkey 8\n\n# Enable writing of statistics records.\n#statistics clockstats cryptostats loopstats peerstats\n","data":{"ntp_servers":["ntp.nict.jp","ntp.jst.mfeed.ad.jp","ntp.ring.gr.jp"]},"delta":"0:00:00.00016","end":"2013-09-20 01:00:10.03129","failed":false,"mode":"644","module":"template","path":"/tmp/foo/bar/baz.conf","skipped":true,"start":"2013-09-20 01:00:10.03112"}
  {"changed":false,"comment":"","delta":"0:00:00.00000","end":"2013-09-20 01:00:10.03131","failed":false,"group":"","groups":null,"home_dir":"","module":"user","name":"foo","shell":"","skipped":true,"start":"2013-09-20 01:00:10.03131","state":"present","system":false,"u.Appends":false,"uid":-1}
  {"changed":false,"delta":"0:00:00.00000","end":"2013-09-20 01:00:10.03131","failed":false,"gid":-1,"module":"group","name":"bar","skipped":true,"start":"2013-09-20 01:00:10.03131","state":"present","system":false}
  {"changed":false,"cmd":"rpm -q ntp","delta":"-2270825:16:28.90268","end":"0001-01-01 00:00:00.00000","failed":false,"module":"yum.installed","name":"ntp","rc":0,"skipped":false,"start":"2013-09-20 01:00:10.03133","stderr":"","stdout":"ntp-4.2.4p8-3.el6.centos.x86_64\n"}
  {"changed":false,"delta":"0:00:00.00000","end":"2013-09-20 01:00:10.04275","failed":false,"module":"yum","name":"ntp","skipped":true,"start":"2013-09-20 01:00:10.04275"}
  {"changed":false,"cmd":"service ntpd status","delta":"0:00:00.01407","end":"2013-09-20 01:00:10.05690","failed":false,"module":"service.state","name":"ntpd","rc":0,"skipped":false,"start":"2013-09-20 01:00:10.04283","stderr":"","stdout":"ntpd (pid  14782) is running...\n"}
  {"changed":false,"delta":"0:00:00.00001","end":"2013-09-20 01:00:10.05703","failed":false,"module":"service.change_state","name":"ntpd","skipped":true,"start":"2013-09-20 01:00:10.05702","state":"started"}
  {"changed":false,"cmd":"chkconfig ntpd --list","delta":"0:00:00.00224","end":"2013-09-20 01:00:10.05932","failed":false,"module":"service.auto_start","rc":0,"skipped":false,"start":"2013-09-20 01:00:10.05708","stderr":"","stdout":"ntpd           \u00090:off\u00091:off\u00092:on\u00093:on\u00094:on\u00095:on\u00096:off\n"}
  {"changed":false,"delta":"0:00:00.00000","end":"2013-09-20 01:00:10.05942","failed":false,"module":"service.change_auto_start","name":"ntpd","skipped":true,"start":"2013-09-20 01:00:10.05942"}
