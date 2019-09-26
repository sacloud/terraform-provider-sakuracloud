// +build linux

package server

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sacloud/libsacloud/sacloud"
	"github.com/skratchdot/open-golang/open"
)

// StartDefaultVNCClient starts OS's default VNC client
func StartDefaultVNCClient(vncProxyInfo *sacloud.VNCProxyResponse) error {

	uri := ""

	for uri == "" {
		// create .vnc tmp-file
		f, err := ioutil.TempFile("", "libsacloud_open_vnc")
		if err != nil {
			return err
		}
		defer f.Close()
		uri = fmt.Sprintf("%s.vnc", f.Name())
		if _, err := os.Stat(uri); err == nil {
			uri = ""
		}
	}
	host := vncProxyInfo.ActualHost()
	body := fmt.Sprintf(vncFileFormat,
		host,
		vncProxyInfo.Port,
		vncProxyInfo.Password,
	)

	ioutil.WriteFile(uri, []byte(body), 0755)
	return open.Start(uri)
}

var vncFileFormat = `[Connection]
Host=%s
Port=%s
Password=%s
`
