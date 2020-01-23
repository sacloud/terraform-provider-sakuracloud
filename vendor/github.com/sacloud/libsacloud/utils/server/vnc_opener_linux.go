// Copyright 2016-2020 The Libsacloud Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
