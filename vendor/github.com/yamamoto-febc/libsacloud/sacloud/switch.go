package sacloud

import (
	"fmt"
	"net"
	"time"
)

// Switch type of switch
type Switch struct {
	*Resource
	Name           string  `json:",omitempty"`
	Description    string  `json:",omitempty"`
	ServerCount    int     `json:",omitempty"`
	ApplianceCount int     `json:",omitempty"`
	Scope          EScope  `json:",omitempty"`
	Subnet         *Subnet `json:",omitempty"`
	UserSubnet     *Subnet `json:",omitempty"`
	//HybridConnection
	ServerClass string     `json:",omitempty"`
	CreatedAt   *time.Time `json:",omitempty"`
	Icon        *Icon      `json:",omitempty"`
	Tags        []string   `json:",omitempty"`
	Subnets     []Subnet   `json:",omitempty"`
	IPv6Nets    []IPv6Net  `json:",omitempty"`
	Internet    *Internet  `json:",omitempty"`
	Bridge      *Bridge    `json:",omitempty"`
}

// Subnet type of Subnet
type Subnet struct {
	*NumberResource
	NetworkAddress string `json:",omitempty"`
	NetworkMaskLen int    `json:",omitempty"`
	DefaultRoute   string `json:",omitempty"`
	//NextHop ???
	//StaticRoute ???
	ServiceClass string `json:",omitempty"`
	IPAddresses  struct {
		Min string `json:",omitempty"`
		Max string `json:",omitempty"`
	}
	Internet *Internet `json:",omitempty"`
}

type IPv6Net struct {
	*NumberResource
	IPv6Prefix    string `json:",omitempty"`
	IPv6PrefixLen int    `json:",omitempty"`
	Scope         string `json:",omitempty"`
	ServiceClass  string `json:",omitempty"`
}

func (s *Switch) GetDefaultIPAddressesForVPCRouter() (string, string, string, error) {

	if s.Subnets == nil || len(s.Subnets) < 1 {
		return "", "", "", fmt.Errorf("switch[%s].Subnets is nil", s.ID)
	}

	baseAddress := net.ParseIP(s.Subnets[0].IPAddresses.Min).To4()
	address1 := net.IPv4(baseAddress[0], baseAddress[1], baseAddress[2], baseAddress[3]+1)
	address2 := net.IPv4(baseAddress[0], baseAddress[1], baseAddress[2], baseAddress[3]+2)

	return baseAddress.String(), address1.String(), address2.String(), nil
}
