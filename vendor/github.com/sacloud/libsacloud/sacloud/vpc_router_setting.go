package sacloud

import (
	"fmt"
	"reflect"
)

// VPCRouterSetting VPCルーター設定
type VPCRouterSetting struct {
	Interfaces         []*VPCRouterInterface        `json:",omitempty"` // NIC設定
	StaticNAT          *VPCRouterStaticNAT          `json:",omitempty"` // スタティックNAT設定
	PortForwarding     *VPCRouterPortForwarding     `json:",omitempty"` // ポートフォワーディング設定
	Firewall           *VPCRouterFirewall           `json:",omitempty"` // ファイアウォール設定
	DHCPServer         *VPCRouterDHCPServer         `json:",omitempty"` // DHCPサーバー設定
	DHCPStaticMapping  *VPCRouterDHCPStaticMapping  `json:",omitempty"` // DHCPスタティックマッピング設定
	L2TPIPsecServer    *VPCRouterL2TPIPsecServer    `json:",omitempty"` // L2TP/IPSecサーバー設定
	PPTPServer         *VPCRouterPPTPServer         `json:",omitempty"` // PPTPサーバー設定
	RemoteAccessUsers  *VPCRouterRemoteAccessUsers  `json:",omitempty"` // リモートアクセスユーザー設定
	SiteToSiteIPsecVPN *VPCRouterSiteToSiteIPsecVPN `json:",omitempty"` // サイト間VPN設定
	StaticRoutes       *VPCRouterStaticRoutes       `json:",omitempty"` // スタティックルート設定
	VRID               *int                         `json:",omitempty"` // VRID
	SyslogHost         string                       `json:",omitempty"` // syslog転送先ホスト

}

// VPCRouterInterface NIC設定
type VPCRouterInterface struct {
	IPAddress        []string `json:",omitempty"` // IPアドレスリスト
	NetworkMaskLen   int      `json:",omitempty"` // ネットワークマスク長
	VirtualIPAddress string   `json:",omitempty"` // 仮想IPアドレス
	IPAliases        []string `json:",omitempty"` // IPエイリアス
}

// AddInterface NIC追加
func (s *VPCRouterSetting) AddInterface(vip string, ipaddress []string, maskLen int) {
	if s.Interfaces == nil {
		s.Interfaces = []*VPCRouterInterface{nil}
	}
	s.Interfaces = append(s.Interfaces, &VPCRouterInterface{
		VirtualIPAddress: vip,
		IPAddress:        ipaddress,
		NetworkMaskLen:   maskLen,
	})
}

// VPCRouterStaticNAT スタティックNAT設定
type VPCRouterStaticNAT struct {
	Config  []*VPCRouterStaticNATConfig `json:",omitempty"` // スタティックNAT設定
	Enabled string                      `json:",omitempty"` // 有効/無効
}

// VPCRouterStaticNATConfig スタティックNAT設定
type VPCRouterStaticNATConfig struct {
	GlobalAddress  string `json:",omitempty"` // グローバルIPアドレス
	PrivateAddress string `json:",omitempty"` // プライベートIPアドレス
}

// AddStaticNAT スタティックNAT設定 追加
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

// RemoveStaticNAT スタティックNAT設定 削除
func (s *VPCRouterSetting) RemoveStaticNAT(globalAddress string, privateAddress string) {
	if s.StaticNAT == nil {
		return
	}

	if s.StaticNAT.Config == nil {
		s.StaticNAT.Enabled = "False"
		return
	}

	dest := []*VPCRouterStaticNATConfig{}
	for _, c := range s.StaticNAT.Config {
		if c.GlobalAddress != globalAddress || c.PrivateAddress != privateAddress {
			dest = append(dest, c)
		}
	}
	s.StaticNAT.Config = dest
	if len(s.StaticNAT.Config) == 0 {
		s.StaticNAT.Enabled = "False"
		s.StaticNAT.Config = nil
		return
	}
	s.StaticNAT.Enabled = "True"
}

// FindStaticNAT スタティックNAT設定検索
func (s *VPCRouterSetting) FindStaticNAT(globalAddress string, privateAddress string) *VPCRouterStaticNATConfig {
	for _, c := range s.StaticNAT.Config {
		if c.GlobalAddress == globalAddress && c.PrivateAddress == privateAddress {
			return c
		}
	}
	return nil
}

// VPCRouterPortForwarding ポートフォワーディング設定
type VPCRouterPortForwarding struct {
	Config  []*VPCRouterPortForwardingConfig `json:",omitempty"` // ポートフォワーディング設定
	Enabled string                           `json:",omitempty"` // 有効/無効
}

// VPCRouterPortForwardingConfig ポートフォワーディング設定
type VPCRouterPortForwardingConfig struct {
	Protocol       string `json:",omitempty"` // プロトコル(tcp/udp)
	GlobalPort     string `json:",omitempty"` // グローバル側ポート
	PrivateAddress string `json:",omitempty"` // プライベートIPアドレス
	PrivatePort    string `json:",omitempty"` // プライベート側ポート
}

// AddPortForwarding ポートフォワーディング 追加
func (s *VPCRouterSetting) AddPortForwarding(protocol string, globalPort string, privateAddress string, privatePort string) {
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
		GlobalPort:     globalPort,
		PrivateAddress: privateAddress,
		PrivatePort:    privatePort,
	})
}

// RemovePortForwarding ポートフォワーディング 削除
func (s *VPCRouterSetting) RemovePortForwarding(protocol string, globalPort string, privateAddress string, privatePort string) {
	if s.PortForwarding == nil {
		return
	}

	if s.PortForwarding.Config == nil {
		s.PortForwarding.Enabled = "False"
		return
	}

	dest := []*VPCRouterPortForwardingConfig{}
	for _, c := range s.PortForwarding.Config {
		if c.Protocol != protocol || c.GlobalPort != globalPort ||
			c.PrivateAddress != privateAddress || c.PrivatePort != privatePort {
			dest = append(dest, c)
		}
	}
	s.PortForwarding.Config = dest
	if len(s.PortForwarding.Config) == 0 {
		s.PortForwarding.Enabled = "False"
		s.PortForwarding.Config = nil
		return
	}
	s.PortForwarding.Enabled = "True"
}

// FindPortForwarding ポートフォワーディング検索
func (s *VPCRouterSetting) FindPortForwarding(protocol string, globalPort string, privateAddress string, privatePort string) *VPCRouterPortForwardingConfig {
	for _, c := range s.PortForwarding.Config {
		if c.Protocol == protocol && c.GlobalPort == globalPort &&
			c.PrivateAddress == privateAddress && c.PrivatePort == privatePort {
			return c
		}
	}
	return nil
}

// VPCRouterFirewall ファイアウォール設定
type VPCRouterFirewall struct {
	Config  []*VPCRouterFirewallSetting `json:",omitempty"` // ファイアウォール設定
	Enabled string                      `json:",omitempty"` // 有効/無効
}

// VPCRouterFirewallSetting ファイアウォール設定
type VPCRouterFirewallSetting struct {
	Receive []*VPCRouterFirewallRule `json:",omitempty"` // 受信ルール
	Send    []*VPCRouterFirewallRule `json:",omitempty"` // 送信ルール
}

// VPCRouterFirewallRule ファイアウォール ルール
type VPCRouterFirewallRule struct {
	Action             string `json:",omitempty"` // 許可/拒否
	Protocol           string `json:",omitempty"` // プロトコル
	SourceNetwork      string `json:",omitempty"` // 送信元ネットワーク
	SourcePort         string `json:",omitempty"` // 送信元ポート
	DestinationNetwork string `json:",omitempty"` // 宛先ネットワーク
	DestinationPort    string `json:",omitempty"` // 宛先ポート
}

func (s *VPCRouterSetting) addFirewallRule(direction string, rule *VPCRouterFirewallRule) {
	if s.Firewall == nil {
		s.Firewall = &VPCRouterFirewall{
			Enabled: "True",
		}
	}
	if s.Firewall.Config == nil || len(s.Firewall.Config) == 0 {
		s.Firewall.Config = []*VPCRouterFirewallSetting{
			{
				Receive: []*VPCRouterFirewallRule{},
				Send:    []*VPCRouterFirewallRule{},
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

func (s *VPCRouterSetting) removeFirewallRule(direction string, rule *VPCRouterFirewallRule) {

	if s.Firewall == nil {
		return
	}

	if s.Firewall.Config == nil {
		s.Firewall.Enabled = "False"
		return
	}

	switch direction {
	case "send":
		dest := []*VPCRouterFirewallRule{}
		for _, c := range s.Firewall.Config[0].Send {
			if !s.isSameRule(rule, c) {
				dest = append(dest, c)
			}
		}
		s.Firewall.Config[0].Send = dest
	case "receive":
		dest := []*VPCRouterFirewallRule{}
		for _, c := range s.Firewall.Config[0].Receive {
			if !s.isSameRule(rule, c) {
				dest = append(dest, c)
			}
		}
		s.Firewall.Config[0].Receive = dest
	}

	if len(s.Firewall.Config) == 0 {
		s.Firewall.Enabled = "False"
		s.Firewall.Config = nil
		return
	}

	if len(s.Firewall.Config[0].Send) == 0 && len(s.Firewall.Config[0].Send) == 0 {
		s.Firewall.Enabled = "False"
		s.Firewall.Config = nil
		return
	}

	s.PortForwarding.Enabled = "True"

}

func (s *VPCRouterSetting) findFirewallRule(direction string, rule *VPCRouterFirewallRule) *VPCRouterFirewallRule {
	switch direction {
	case "send":
		for _, c := range s.Firewall.Config[0].Send {
			if s.isSameRule(rule, c) {
				return c
			}
		}
	case "receive":
		for _, c := range s.Firewall.Config[0].Receive {
			if s.isSameRule(rule, c) {
				return c
			}
		}
	}

	return nil

}

func (s *VPCRouterSetting) isSameRule(r1 *VPCRouterFirewallRule, r2 *VPCRouterFirewallRule) bool {
	return r1.Action == r2.Action &&
		r1.Protocol == r2.Protocol &&
		r1.SourceNetwork == r2.SourceNetwork &&
		r1.SourcePort == r2.SourcePort &&
		r1.DestinationNetwork == r2.DestinationNetwork &&
		r1.DestinationPort == r2.DestinationPort
}

// AddFirewallRuleSend 送信ルール 追加
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

// RemoveFirewallRuleSend 送信ルール 削除
func (s *VPCRouterSetting) RemoveFirewallRuleSend(isAllow bool, protocol string, sourceNetwork string, sourcePort string, destNetwork string, destPort string) {
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

	s.removeFirewallRule("send", rule)
}

// FindFirewallRuleSend 送信ルール 検索
func (s *VPCRouterSetting) FindFirewallRuleSend(isAllow bool, protocol string, sourceNetwork string, sourcePort string, destNetwork string, destPort string) *VPCRouterFirewallRule {
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

	return s.findFirewallRule("send", rule)
}

// AddFirewallRuleReceive 受信ルール 追加
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

// RemoveFirewallRuleReceive 受信ルール 削除
func (s *VPCRouterSetting) RemoveFirewallRuleReceive(isAllow bool, protocol string, sourceNetwork string, sourcePort string, destNetwork string, destPort string) {
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

	s.removeFirewallRule("receive", rule)
}

// FindFirewallRuleReceive 受信ルール 検索
func (s *VPCRouterSetting) FindFirewallRuleReceive(isAllow bool, protocol string, sourceNetwork string, sourcePort string, destNetwork string, destPort string) *VPCRouterFirewallRule {
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

	return s.findFirewallRule("receive", rule)
}

// VPCRouterDHCPServer DHCPサーバー設定
type VPCRouterDHCPServer struct {
	Config  []*VPCRouterDHCPServerConfig `json:",omitempty"` // DHCPサーバー設定
	Enabled string                       `json:",omitempty"` // 有効/無効
}

// VPCRouterDHCPServerConfig DHCPサーバー設定
type VPCRouterDHCPServerConfig struct {
	Interface  string `json:",omitempty"` // 対象NIC
	RangeStart string `json:",omitempty"` // 割り当て範囲 開始アドレス
	RangeStop  string `json:",omitempty"` // 割り当て範囲 終了アドレス
}

// AddDHCPServer DHCPサーバー設定追加
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

// RemoveDHCPServer DHCPサーバー設定削除
func (s *VPCRouterSetting) RemoveDHCPServer(nicIndex int, rangeStart string, rangeStop string) {
	if s.DHCPServer == nil {
		return
	}

	if s.DHCPServer.Config == nil {
		s.DHCPServer.Enabled = "False"
		return
	}

	dest := []*VPCRouterDHCPServerConfig{}
	for _, c := range s.DHCPServer.Config {
		if c.Interface != fmt.Sprintf("eth%d", nicIndex) || c.RangeStart != rangeStart || c.RangeStop != rangeStop {
			dest = append(dest, c)
		}
	}
	s.DHCPServer.Config = dest

	if len(s.DHCPServer.Config) == 0 {
		s.DHCPServer.Enabled = "False"
		s.DHCPServer.Config = nil
		return
	}
	s.DHCPServer.Enabled = "True"

}

// FindDHCPServer DHCPサーバー設定 検索
func (s *VPCRouterSetting) FindDHCPServer(nicIndex int, rangeStart string, rangeStop string) *VPCRouterDHCPServerConfig {
	for _, c := range s.DHCPServer.Config {
		if c.Interface == fmt.Sprintf("eth%d", nicIndex) && c.RangeStart == rangeStart && c.RangeStop == rangeStop {
			return c
		}
	}
	return nil
}

// VPCRouterDHCPStaticMapping DHCPスタティックマッピング設定
type VPCRouterDHCPStaticMapping struct {
	Config  []*VPCRouterDHCPStaticMappingConfig `json:",omitempty"` // DHCPスタティックマッピング設定
	Enabled string                              `json:",omitempty"` // 有効/無効
}

// VPCRouterDHCPStaticMappingConfig DHCPスタティックマッピング設定
type VPCRouterDHCPStaticMappingConfig struct {
	IPAddress  string `json:",omitempty"` // 割り当てIPアドレス
	MACAddress string `json:",omitempty"` // ソースMACアドレス
}

// AddDHCPStaticMapping DHCPスタティックマッピング設定追加
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

// RemoveDHCPStaticMapping DHCPスタティックマッピング設定 削除
func (s *VPCRouterSetting) RemoveDHCPStaticMapping(ipAddress string, macAddress string) {
	if s.DHCPStaticMapping == nil {
		return
	}

	if s.DHCPStaticMapping.Config == nil {
		s.DHCPStaticMapping.Enabled = "False"
		return
	}

	dest := []*VPCRouterDHCPStaticMappingConfig{}
	for _, c := range s.DHCPStaticMapping.Config {
		if c.IPAddress != ipAddress || c.MACAddress != macAddress {
			dest = append(dest, c)
		}
	}
	s.DHCPStaticMapping.Config = dest

	if len(s.DHCPStaticMapping.Config) == 0 {
		s.DHCPStaticMapping.Enabled = "False"
		s.DHCPStaticMapping.Config = nil
		return
	}
	s.DHCPStaticMapping.Enabled = "True"

}

// FindDHCPStaticMapping DHCPスタティックマッピング設定 検索
func (s *VPCRouterSetting) FindDHCPStaticMapping(ipAddress string, macAddress string) *VPCRouterDHCPStaticMappingConfig {
	for _, c := range s.DHCPStaticMapping.Config {
		if c.IPAddress == ipAddress && c.MACAddress == macAddress {
			return c
		}
	}
	return nil
}

// VPCRouterL2TPIPsecServer L2TP/IPSecサーバー設定
type VPCRouterL2TPIPsecServer struct {
	Config  *VPCRouterL2TPIPsecServerConfig `json:",omitempty"` // L2TP/IPSecサーバー設定
	Enabled string                          `json:",omitempty"` // 有効/無効
}

// VPCRouterL2TPIPsecServerConfig L2TP/IPSecサーバー設定
type VPCRouterL2TPIPsecServerConfig struct {
	PreSharedSecret string `json:",omitempty"` // 事前共有シークレット
	RangeStart      string `json:",omitempty"` // 割り当て範囲 開始IPアドレス
	RangeStop       string `json:",omitempty"` // 割り当て範囲 終了IPアドレス
}

// EnableL2TPIPsecServer L2TP/IPSecサーバー設定 有効化
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

// DisableL2TPIPsecServer L2TP/IPSecサーバー設定 無効化
func (s *VPCRouterSetting) DisableL2TPIPsecServer() {
	if s.L2TPIPsecServer == nil {
		s.L2TPIPsecServer = &VPCRouterL2TPIPsecServer{}
	}
	s.L2TPIPsecServer.Enabled = "False"
	s.L2TPIPsecServer.Config = nil
}

// VPCRouterPPTPServer PPTPサーバー設定
type VPCRouterPPTPServer struct {
	Config  *VPCRouterPPTPServerConfig `json:",omitempty"` // PPTPサーバー設定
	Enabled string                     `json:",omitempty"` // 有効/無効
}

// VPCRouterPPTPServerConfig PPTPサーバー設定
type VPCRouterPPTPServerConfig struct {
	RangeStart string `json:",omitempty"` // 割り当て範囲 開始IPアドレス
	RangeStop  string `json:",omitempty"` // 割り当て範囲 終了IPアドレス
}

// EnablePPTPServer PPTPサーバー設定 有効化
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

// DisablePPTPServer PPTPサーバー設定 無効化
func (s *VPCRouterSetting) DisablePPTPServer() {
	if s.PPTPServer == nil {
		s.PPTPServer = &VPCRouterPPTPServer{}
	}
	s.PPTPServer.Enabled = "False"
	s.PPTPServer.Config = nil
}

// VPCRouterRemoteAccessUsers リモートアクセスユーザー設定
type VPCRouterRemoteAccessUsers struct {
	Config  []*VPCRouterRemoteAccessUsersConfig `json:",omitempty"` // リモートアクセスユーザー設定
	Enabled string                              `json:",omitempty"` // 有効/無効
}

// VPCRouterRemoteAccessUsersConfig リモートアクセスユーザー設定
type VPCRouterRemoteAccessUsersConfig struct {
	UserName string `json:",omitempty"` // ユーザー名
	Password string `json:",omitempty"` // パスワード
}

// AddRemoteAccessUser リモートアクセスユーザー設定 追加
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

// RemoveRemoteAccessUser リモートアクセスユーザー設定 削除
func (s *VPCRouterSetting) RemoveRemoteAccessUser(userName string, password string) {
	if s.RemoteAccessUsers == nil {
		return
	}

	if s.RemoteAccessUsers.Config == nil {
		s.RemoteAccessUsers.Enabled = "False"
		return
	}

	dest := []*VPCRouterRemoteAccessUsersConfig{}
	for _, c := range s.RemoteAccessUsers.Config {
		if c.UserName != userName || c.Password != password {
			dest = append(dest, c)
		}
	}
	s.RemoteAccessUsers.Config = dest

	if len(s.RemoteAccessUsers.Config) == 0 {
		s.RemoteAccessUsers.Enabled = "False"
		s.RemoteAccessUsers.Config = nil
		return
	}
	s.RemoteAccessUsers.Enabled = "True"
}

// FindRemoteAccessUser リモートアクセスユーザー設定 検索
func (s *VPCRouterSetting) FindRemoteAccessUser(userName string, password string) *VPCRouterRemoteAccessUsersConfig {
	for _, c := range s.RemoteAccessUsers.Config {
		if c.UserName == userName && c.Password == password {
			return c
		}
	}
	return nil
}

// VPCRouterSiteToSiteIPsecVPN サイト間VPC設定
type VPCRouterSiteToSiteIPsecVPN struct {
	Config  []*VPCRouterSiteToSiteIPsecVPNConfig `json:",omitempty"` // サイト間VPC設定
	Enabled string                               `json:",omitempty"` // 有効/無効
}

// VPCRouterSiteToSiteIPsecVPNConfig サイト間VPC設定
type VPCRouterSiteToSiteIPsecVPNConfig struct {
	LocalPrefix     []string `json:",omitempty"` // ローカルプレフィックス リスト
	Peer            string   `json:",omitempty"` // 対向IPアドレス
	PreSharedSecret string   `json:",omitempty"` // 事前共有シークレット
	RemoteID        string   `json:",omitempty"` // 対向ID
	Routes          []string `json:",omitempty"` // 対向プレフィックス リスト
}

// AddSiteToSiteIPsecVPN サイト間VPC設定 追加
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

// RemoveSiteToSiteIPsecVPN サイト間VPC設定 削除
func (s *VPCRouterSetting) RemoveSiteToSiteIPsecVPN(localPrefix []string, peer string, preSharedSecret string, remoteID string, routes []string) {
	config := &VPCRouterSiteToSiteIPsecVPNConfig{
		LocalPrefix:     localPrefix,
		Peer:            peer,
		PreSharedSecret: preSharedSecret,
		RemoteID:        remoteID,
		Routes:          routes,
	}

	if s.SiteToSiteIPsecVPN == nil {
		return
	}

	if s.SiteToSiteIPsecVPN.Config == nil {
		s.SiteToSiteIPsecVPN.Enabled = "False"
		return
	}

	dest := []*VPCRouterSiteToSiteIPsecVPNConfig{}
	for _, c := range s.SiteToSiteIPsecVPN.Config {
		if !s.isSameSiteToSiteIPsecVPNConfig(c, config) {
			dest = append(dest, c)
		}
	}
	s.SiteToSiteIPsecVPN.Config = dest

	if len(s.SiteToSiteIPsecVPN.Config) == 0 {
		s.SiteToSiteIPsecVPN.Enabled = "False"
		s.SiteToSiteIPsecVPN.Config = nil
		return
	}
	s.SiteToSiteIPsecVPN.Enabled = "True"
}

// FindSiteToSiteIPsecVPN  サイト間VPC設定 検索
func (s *VPCRouterSetting) FindSiteToSiteIPsecVPN(localPrefix []string, peer string, preSharedSecret string, remoteID string, routes []string) *VPCRouterSiteToSiteIPsecVPNConfig {
	config := &VPCRouterSiteToSiteIPsecVPNConfig{
		LocalPrefix:     localPrefix,
		Peer:            peer,
		PreSharedSecret: preSharedSecret,
		RemoteID:        remoteID,
		Routes:          routes,
	}

	for _, c := range s.SiteToSiteIPsecVPN.Config {
		if s.isSameSiteToSiteIPsecVPNConfig(c, config) {
			return c
		}
	}
	return nil
}

func (s *VPCRouterSetting) isSameSiteToSiteIPsecVPNConfig(c1 *VPCRouterSiteToSiteIPsecVPNConfig, c2 *VPCRouterSiteToSiteIPsecVPNConfig) bool {
	return reflect.DeepEqual(c1.LocalPrefix, c2.LocalPrefix) &&
		c1.Peer == c2.Peer &&
		c1.PreSharedSecret == c2.PreSharedSecret &&
		c1.RemoteID == c2.RemoteID &&
		reflect.DeepEqual(c1.Routes, c2.Routes)
}

// VPCRouterStaticRoutes スタティックルート設定
type VPCRouterStaticRoutes struct {
	Config  []*VPCRouterStaticRoutesConfig `json:",omitempty"` // スタティックルート設定
	Enabled string                         `json:",omitempty"` // 有効/無効
}

// VPCRouterStaticRoutesConfig スタティックルート設定
type VPCRouterStaticRoutesConfig struct {
	Prefix  string `json:",omitempty"` // プレフィックス
	NextHop string `json:",omitempty"` // ネクストホップ
}

// AddStaticRoute スタティックルート設定 追加
func (s *VPCRouterSetting) AddStaticRoute(prefix string, nextHop string) {
	if s.StaticRoutes == nil {
		s.StaticRoutes = &VPCRouterStaticRoutes{
			Enabled: "True",
		}
	}
	if s.StaticRoutes.Config == nil {
		s.StaticRoutes.Config = []*VPCRouterStaticRoutesConfig{}
	}
	s.StaticRoutes.Config = append(s.StaticRoutes.Config, &VPCRouterStaticRoutesConfig{
		Prefix:  prefix,
		NextHop: nextHop,
	})
}

// RemoveStaticRoute スタティックルート設定 削除
func (s *VPCRouterSetting) RemoveStaticRoute(prefix string, nextHop string) {
	if s.StaticRoutes == nil {
		return
	}

	if s.StaticRoutes.Config == nil {
		s.StaticRoutes.Enabled = "False"
		return
	}

	dest := []*VPCRouterStaticRoutesConfig{}
	for _, c := range s.StaticRoutes.Config {
		if c.Prefix != prefix || c.NextHop != nextHop {
			dest = append(dest, c)
		}
	}
	s.StaticRoutes.Config = dest

	if len(s.StaticRoutes.Config) == 0 {
		s.StaticRoutes.Enabled = "False"
		s.StaticRoutes.Config = nil
		return
	}
	s.StaticRoutes.Enabled = "True"
}

// FindStaticRoute スタティックルート設定 検索
func (s *VPCRouterSetting) FindStaticRoute(prefix string, nextHop string) *VPCRouterStaticRoutesConfig {
	for _, c := range s.StaticRoutes.Config {
		if c.Prefix == prefix && c.NextHop == nextHop {
			return c
		}
	}
	return nil
}
