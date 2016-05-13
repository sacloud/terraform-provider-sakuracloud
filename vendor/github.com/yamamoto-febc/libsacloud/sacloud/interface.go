package sacloud

// Interface type of server nic
type Interface struct {
	*Resource
	MACAddress    string        `json:",omitempty"`
	IPAddress     string        `json:",omitempty"`
	UserIPAddress string        `json:",omitempty"`
	HostName      string        `json:",omitempty"`
	Server        *Server       `json:",omitempty"`
	Switch        *Switch       `json:",omitempty"`
	PacketFilter  *PacketFilter `json:",omitempty"`
}

func (i *Interface) SetNewServerID(id string) {
	i.Server = &Server{Resource: &Resource{ID: id}}
}

func (i *Interface) SetNewSwitchID(id string) {
	i.Switch = &Switch{Resource: &Resource{ID: id}}
}
