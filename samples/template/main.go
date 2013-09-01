package main

import (
	"fmt"
	"io/ioutil"
	"text/template"
	"github.com/hnakamur/commango/templateutil"
	"github.com/hnakamur/commango/hashutil"
)

func main() {
	tmpl, err := template.ParseFiles("sample1.tmpl")
	if err != nil {
		panic(err)
	}

	data := map[string]interface{} {
		"ntp_servers": []string{
			"ntp.nict.jp",
			"ntp.jst.mfeed.ad.jp",
			"ntp.ring.gr.jp",
		},
	}
	output, err := templateutil.RenderToBytes(tmpl, data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("output: %s", output)
	sum := hashutil.CalcSha256Sum(output)
	ioutil.WriteFile("a.conf", output, 0644)
	fmt.Printf("sum: %s\n", sum)
}
