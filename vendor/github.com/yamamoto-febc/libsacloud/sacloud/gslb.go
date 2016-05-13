package sacloud

import "time"

// GSLB type of GSLB(CommonServiceItem)
type GSLB struct {
	*Resource
	Name         string
	Description  string       `json:",omitempty"`
	Status       GSLBStatus   `json:",omitempty"`
	Provider     GSLBProvider `json:",omitempty"`
	Settings     GSLBSettings `json:",omitempty"`
	ServiceClass string       `json:",omitempty"`
	CreatedAt    *time.Time   `json:",omitempty"`
	ModifiedAt   *time.Time   `json:",omitempty"`
	Icon         *Icon        `json:",omitempty"`
	Tags         []string     `json:",omitempty"`
}

// GSLBSettings type of GSLBSettings
type GSLBSettings struct {
	GSLB GSLBRecordSets `json:",omitempty"`
}

// GSLBStatus type of GSLBStatus
type GSLBStatus struct {
	FQDN string `json:",omitempty"`
}

// GSLBProvider type of GSLBProvider
type GSLBProvider struct {
	Class string `json:",omitempty"`
}

// CreateNewGSLB Create new GLSB(CommonServiceItem)
func CreateNewGSLB(gslbName string) *GSLB {
	return &GSLB{
		Resource: &Resource{ID: ""},
		Name:     gslbName,
		Provider: GSLBProvider{
			Class: "gslb",
		},
		Settings: GSLBSettings{
			GSLB: GSLBRecordSets{
				DelayLoop:   10,
				HealthCheck: defaultGSLBHealthCheck,
				Weighted:    "True",
				Servers:     []GSLBServer{},
			},
		},
	}

}

func AllowGSLBHealthCheckProtocol() []string {
	return []string{"http", "https", "ping", "tcp"}
}

// HasGSLBServer return has server
func (d *GSLB) HasGSLBServer() bool {
	return len(d.Settings.GSLB.Servers) > 0
}

func (d *GSLB) CreateGSLBServer(ip string) *GSLBServer {
	return &GSLBServer{
		IPAddress: ip,
		Enabled:   "True",
		Weight:    "1",
	}
}

func (d *GSLB) AddGSLBServer(server *GSLBServer) {
	var isExist = false
	for i := range d.Settings.GSLB.Servers {
		if d.Settings.GSLB.Servers[i].IPAddress == server.IPAddress {
			d.Settings.GSLB.Servers[i].Enabled = server.Enabled
			d.Settings.GSLB.Servers[i].Weight = server.Weight
			isExist = true
		}
	}

	if !isExist {
		d.Settings.GSLB.Servers = append(d.Settings.GSLB.Servers, *server)
	}
}

func (d *GSLB) ClearGSLBServer() {
	d.Settings.GSLB.Servers = []GSLBServer{}
}

// GSLBRecordSets type of GSLBRecordSets
type GSLBRecordSets struct {
	DelayLoop   int             `json:",omitempty"`
	HealthCheck GSLBHealthCheck `json:",omitempty"`
	Weighted    string          `json:",omitempty"`
	Servers     []GSLBServer    `json:",omitempty"`
}

// AddServer Add server to GSLB
func (g *GSLBRecordSets) AddServer(ip string) {
	var record GSLBServer
	var isExist = false
	for i := range g.Servers {
		if g.Servers[i].IPAddress == ip {
			isExist = true
		}
	}

	if !isExist {
		record = GSLBServer{
			IPAddress: ip,
			Enabled:   "True",
			Weight:    "1",
		}
		g.Servers = append(g.Servers, record)
	}
}

//DeleteServer Delete server from GSLB
func (g *GSLBRecordSets) DeleteServer(ip string) {
	res := []GSLBServer{}
	for i := range g.Servers {
		if g.Servers[i].IPAddress != ip {
			res = append(res, g.Servers[i])
		}
	}

	g.Servers = res
}

// GSLBServer type of GSLBServer
type GSLBServer struct {
	IPAddress string `json:",omitempty"`
	Enabled   string `json:",omitempty"`
	Weight    string `json:",omitempty"`
}

// GSLBHealthCheck type of GSLBHealthCheck
type GSLBHealthCheck struct {
	Protocol string `json:",omitempty"`
	Host     string `json:",omitempty"`
	Path     string `json:",omitempty"`
	Status   string `json:",omitempty"`
	Port     string `json:",omitempty"`
}

var defaultGSLBHealthCheck = GSLBHealthCheck{
	Protocol: "http",
	Host:     "",
	Path:     "/",
	Status:   "200",
}
