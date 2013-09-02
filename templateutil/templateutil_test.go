package templateutil_test

import (
	"io/ioutil"
	"os"
	"testing"
	"github.com/hnakamur/commango/templateutil"
)

const SAMPLE1_TMPL =
`{{range .ntp_servers}}{{/*
*/}}server {{.}}
{{end}}`

func TestWriteIfChanged(t *testing.T) {
	tmpl, err := templateutil.NewWithString("sample1.tmpl", SAMPLE1_TMPL)
	if err != nil {
		t.Fatal(err)
	}

	f, err := ioutil.TempFile("", "output-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	data := map[string]interface{} {
		"ntp_servers": []string{
			"ntp.nict.jp",
			"ntp.jst.mfeed.ad.jp",
			"ntp.ring.gr.jp",
		},
	}

	changed, err := templateutil.WriteIfChanged(tmpl, data, f.Name(), 0644)
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("First write must change the file content")
	}

	changed, err = templateutil.WriteIfChanged(tmpl, data, f.Name(), 0644)
	if err != nil {
		t.Fatal(err)
	}
	if changed {
		t.Fatal("Second write must not overwrite because the content is the same.")
	}
}
