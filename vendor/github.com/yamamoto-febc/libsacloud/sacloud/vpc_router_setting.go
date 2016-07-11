package sacloud

import "fmt"

type VPCRouterSetting struct {
	Interfaces         []*VPCRouterInterface        `json:",omitempty"`
	StaticNAT          *VPCRouterStaticNAT          `json:",omitempty"`
	PortForwarding     *VPCRouterPortForwarding     `json:",omitempty"`
	Firewall           *VPCRouterFirewall           `json:",omitempty"`
	DHCPServer         *VPCRouterDHCPServer         `json:",omitempty"`
	DHCPStaticMapping  *VPCRouterDHCPStaticMapping  `json:",omitempty"`
	L2TPIPsecServer    *VPCRouterL2TPIPsecServer    `json:",omitempty"`
	PPTPServer         *VPCRouterPPTPServer         `json:",omitempty"`
	RemoteAccessUsers  *VPCRouterRemoteAccessUsers  `json:",omitempty"`
	SiteToSiteIPsecVPN *VPCRouterSiteToSiteIPsecVPN `json:",omitempty"`
	VRID               *int                         `json:",omitempty"`
}

type VPCRouterInterface struct {
	IPAddress        []string `json:",omitempty"`
	NetworkMaskLen   int      `json:",omitempty"`
	VirtualIPAddress string   `json:",omitempty"`
}

type VPCRouterStaticNAT struct {
	Config  []*VPCRouterStaticNATConfig `json:",omitempty"`
	Enabled string                      `json:",omitempty"`
}
type VPCRouterStaticNATConfig struct {
	GlobalAddress  string `json:",omitempty"`
	PrivateAddress string `json:",omitempty"`
}

func (s *VPCRouterSetting) AddStaticNAT(globalAddress string, privateAddress string) {
	if s.StaticNAT == nil {
		s.StaticNAT = &VPCRouterStaticNAT{
			Enabled: "True",
		}
	}

	if s.StaticNAT.Config == nil {
		s.StaticNAT.Config = []*VPCRouterStaticNATConfig{}
	}

	s.StaticNAT.Config = append(s.StaticNAT.Config, &VPCRouterStaticNATConfig{
		GlobalAddress:  globalAddress,
		PrivateAddress: privateAddress,
	})
}

type VPCRouterPortForwarding struct {
	Config  []*VPCRouterPortForwardingConfig `json:",omitempty"`
	Enabled string                           `json:",omitempty"`
}
type VPCRouterPortForwardingConfig struct {
	Protocol       string `json:",omitempty"` // tcp/udp only
	GlobalPort     string `json:",omitempty"`
	PrivateAddress string `json:",omitempty"`
	PrivatePort    string `json:",omitempty"`
}

func (s *VPCRouterSetting) AddPortForwarding(protocol string, globalPort int, privateAddress string, privatePort int) {
	if s.PortForwarding == nil {
		s.PortForwarding = &VPCRouterPortForwarding{
			Enabled: "True",
		}
	}

	if s.PortForwarding.Config == nil {
		s.PortForwarding.Config = []*VPCRouterPortForwardingConfig{}
	}

	s.PortForwarding.Config = append(s.PortForwarding.Config, &VPCRouterPortForwardingConfig{
		Protocol:       protocol,
		GlobalPort:     fmt.Sprintf("%d", globalPort),
		PrivateAddress: privateAddress,
		PrivatePort:    fmt.Sprintf("%d", privatePort),
	})
}

type VPCRouterFirewall struct {
	Config  []*VPCRouterFirewallSetting `json:",omitempty"`
	Enabled string                      `json:",omitempty"`
}
type VPCRouterFirewallSetting struct {
	Receive []*VPCRouterFirewallRule `json:",omitempty"`
	Send    []*VPCRouterFirewallRule `json:",omitempty"`
}
type VPCRouterFirewallRule struct {
	Action             string `json:",omitempty"`
	Protocol           string `json:",omitempty"`
	SourceNetwork      string `json:",omitempty"`
	SourcePort         string `json:",omitempty"`
	DestinationNetwork string `json:",omitempty"`
	DestinationPort    string `json:",omitempty"`
}

func (s *VPCRouterSetting) addFirewallRule(direction string, rule *VPCRouterFirewallRule) {
	if s.Firewall == nil {
		s.Firewall = &VPCRouterFirewall{
			Enabled: "True",
			Config: []*VPCRouterFirewallSetting{
				{
					Receive: []*VPCRouterFirewallRule{},
					Send:    []*VPCRouterFirewallRule{},
				},
			},
		}
	}
	switch direction {
	case "send":
		s.Firewall.Config[0].Send = append(s.Firewall.Config[0].Send, rule)
	case "receive":
		s.Firewall.Config[0].Receive = append(s.Firewall.Config[0].Receive, rule)
	}
}

func (s *VPCRouterSetting) AddFirewallRuleSend(isAllow bool, protocol string, sourceNetwork string, sourcePort string, destNetwork string, destPort string) {
	action := "deny"
	if isAllow {
		action = "allow"
	}
	rule := &VPCRouterFirewallRule{
		Action:             action,
		Protocol:           protocol,
		SourceNetwork:      sourceNetwork,
		SourcePort:         sourcePort,
		DestinationNetwork: destNetwork,
		DestinationPort:    destPort,
	}

	s.addFirewallRule("send", rule)
}
func (s *VPCRouterSetting) AddFirewallRuleReceive(isAllow bool, protocol string, sourceNetwork string, sourcePort string, destNetwork string, destPort string) {
	action := "deny"
	if isAllow {
		action = "allow"
	}
	rule := &VPCRouterFirewallRule{
		Action:             action,
		Protocol:           protocol,
		SourceNetwork:      sourceNetwork,
		SourcePort:         sourcePort,
		DestinationNetwork: destNetwork,
		DestinationPort:    destPort,
	}

	s.addFirewallRule("receive", rule)
}

type VPCRouterDHCPServer struct {
	Config  []*VPCRouterDHCPServerConfig `json:",omitempty"`
	Enabled string                       `json:",omitempty"`
}
type VPCRouterDHCPServerConfig struct {
	Interface  string `json:",omitempty"`
	RangeStart string `json:",omitempty"`
	RangeStop  string `json:",omitempty"`
}

func (s *VPCRouterSetting) AddDHCPServer(nicIndex int, rangeStart string, rangeStop string) {
	if s.DHCPServer == nil {
		s.DHCPServer = &VPCRouterDHCPServer{
			Enabled: "True",
		}
	}
	if s.DHCPServer.Config == nil {
		s.DHCPServer.Config = []*VPCRouterDHCPServerConfig{}
	}

	nic := fmt.Sprintf("eth%d", nicIndex)
	s.DHCPServer.Config = append(s.DHCPServer.Config, &VPCRouterDHCPServerConfig{
		Interface:  nic,
		RangeStart: rangeStart,
		RangeStop:  rangeStop,
	})

}

type VPCRouterDHCPStaticMapping struct {
	Config  []*VPCRouterDHCPStaticMappingConfig `json:",omitempty"`
	Enabled string                              `json:",omitempty"`
}
type VPCRouterDHCPStaticMappingConfig struct {
	IPAddress  string `json:",omitempty"`
	MACAddress string `json:",omitempty"`
}

func (s *VPCRouterSetting) AddDHCPStaticMapping(ipAddress string, macAddress string) {
	if s.DHCPStaticMapping == nil {
		s.DHCPStaticMapping = &VPCRouterDHCPStaticMapping{
			Enabled: "True",
		}
	}
	if s.DHCPStaticMapping.Config == nil {
		s.DHCPStaticMapping.Config = []*VPCRouterDHCPStaticMappingConfig{}
	}

	s.DHCPStaticMapping.Config = append(s.DHCPStaticMapping.Config, &VPCRouterDHCPStaticMappingConfig{
		IPAddress:  ipAddress,
		MACAddress: macAddress,
	})
}

type VPCRouterL2TPIPsecServer struct {
	Config  *VPCRouterL2TPIPsecServerConfig `json:",omitempty"`
	Enabled string                          `json:",omitempty"`
}

type VPCRouterL2TPIPsecServerConfig struct {
	PreSharedSecret string `json:",omitempty"`
	RangeStart      string `json:",omitempty"`
	RangeStop       string `json:",omitempty"`
}

func (s *VPCRouterSetting) EnableL2TPIPsecServer(preSharedSecret string, rangeStart string, rangeStop string) {
	if s.L2TPIPsecServer == nil {
		s.L2TPIPsecServer = &VPCRouterL2TPIPsecServer{
			Enabled: "True",
		}
	}
	s.L2TPIPsecServer.Config = &VPCRouterL2TPIPsecServerConfig{
		PreSharedSecret: preSharedSecret,
		RangeStart:      rangeStart,
		RangeStop:       rangeStop,
	}
}

func (s *VPCRouterSetting) DisableL2TPIPsecServer() {
	if s.L2TPIPsecServer == nil {
		s.L2TPIPsecServer = &VPCRouterL2TPIPsecServer{
			Enabled: "False",
		}
	}
	s.L2TPIPsecServer.Config = nil
}

type VPCRouterPPTPServer struct {
	Config  *VPCRouterPPTPServerConfig `json:",omitempty"`
	Enabled string                     `json:",omitempty"`
}
type VPCRouterPPTPServerConfig struct {
	RangeStart string `json:",omitempty"`
	RangeStop  string `json:",omitempty"`
}

func (s *VPCRouterSetting) EnablePPTPServer(rangeStart string, rangeStop string) {
	if s.PPTPServer == nil {
		s.PPTPServer = &VPCRouterPPTPServer{
			Enabled: "True",
		}
	}
	s.PPTPServer.Config = &VPCRouterPPTPServerConfig{
		RangeStart: rangeStart,
		RangeStop:  rangeStop,
	}
}

type VPCRouterRemoteAccessUsers struct {
	Config  []*VPCRouterRemoteAccessUsersConfig `json:",omitempty"`
	Enabled string                              `json:",omitempty"`
}
type VPCRouterRemoteAccessUsersConfig struct {
	UserName string `json:",omitempty"`
	Password string `json:",omitempty"`
}

func (s *VPCRouterSetting) AddRemoteAccessUser(userName string, password string) {
	if s.RemoteAccessUsers == nil {
		s.RemoteAccessUsers = &VPCRouterRemoteAccessUsers{
			Enabled: "True",
		}
	}
	if s.RemoteAccessUsers.Config == nil {
		s.RemoteAccessUsers.Config = []*VPCRouterRemoteAccessUsersConfig{}
	}
	s.RemoteAccessUsers.Config = append(s.RemoteAccessUsers.Config, &VPCRouterRemoteAccessUsersConfig{
		UserName: userName,
		Password: password,
	})
}

type VPCRouterSiteToSiteIPsecVPN struct {
	Config  []*VPCRouterSiteToSiteIPsecVPNConfig `json:",omitempty"`
	Enabled string                               `json:",omitempty"`
}

type VPCRouterSiteToSiteIPsecVPNConfig struct {
	LocalPrefix     []string `json:",omitempty"`
	Peer            string   `json:",omitempty"`
	PreSharedSecret string   `json:",omitempty"`
	RemoteID        string   `json:",omitempty"`
	Routes          []string `json:",omitempty"`
}

func (s *VPCRouterSetting) AddSiteToSiteIPsecVPN(localPrefix []string, peer string, preSharedSecret string, remoteID string, routes []string) {
	if s.SiteToSiteIPsecVPN == nil {
		s.SiteToSiteIPsecVPN = &VPCRouterSiteToSiteIPsecVPN{
			Enabled: "True",
		}
	}
	if s.SiteToSiteIPsecVPN.Config == nil {
		s.SiteToSiteIPsecVPN.Config = []*VPCRouterSiteToSiteIPsecVPNConfig{}
	}

	s.SiteToSiteIPsecVPN.Config = append(s.SiteToSiteIPsecVPN.Config, &VPCRouterSiteToSiteIPsecVPNConfig{
		LocalPrefix:     localPrefix,
		Peer:            peer,
		PreSharedSecret: preSharedSecret,
		RemoteID:        remoteID,
		Routes:          routes,
	})
}
