// +build !windows,!linux

package server

import (
	"fmt"

	"github.com/sacloud/libsacloud/sacloud"
	"github.com/skratchdot/open-golang/open"
)

// StartDefaultVNCClient starts OS's default VNC client
func StartDefaultVNCClient(vncProxyInfo *sacloud.VNCProxyResponse) error {
	host := vncProxyInfo.ActualHost()
	uri := fmt.Sprintf("vnc://:%s@%s:%s",
		vncProxyInfo.Password,
		host,
		vncProxyInfo.Port)
	return open.Start(uri)
}
