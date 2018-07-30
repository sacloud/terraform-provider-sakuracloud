package server

import (
	"fmt"
	"io/ioutil"

	"github.com/skratchdot/open-golang/open"
)

// RDPOpener information of RDP connection
type RDPOpener struct {
	IPAddress       string
	User            string
	Port            int
	RDPFileTemplate string
}

// RDPFileContent represents .rdp file contents
func (c *RDPOpener) RDPFileContent() string {
	addr := c.IPAddress
	if c.Port > 0 {
		addr = fmt.Sprintf("%s:%d", c.IPAddress, c.Port)
	}

	template := c.RDPFileTemplate
	if template == "" {
		template = defaultRDPTemplate
	}
	return fmt.Sprintf(template, addr, c.User)
}

var defaultRDPTemplate = `
full address:s:%s
username:s:%s
autoreconnection enabled:i:1
`

// StartDefaultClient starts OS's default RDP client
func (c *RDPOpener) StartDefaultClient() error {
	uri := ""

	// create .rdp tmp-file
	f, err := ioutil.TempFile("", "usacloud_open_rdp")
	if err != nil {
		return err
	}
	defer f.Close()

	uri = fmt.Sprintf("%s.rdp", f.Name())
	rdpContent := c.RDPFileContent()

	ioutil.WriteFile(uri, []byte(rdpContent), 0755)
	return open.Start(uri)
}
