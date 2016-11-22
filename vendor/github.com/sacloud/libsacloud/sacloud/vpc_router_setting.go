package sacloud

import (
	"fmt"
	"reflect"
)

// VPCRouterSetting VPCルーター設定
type VPCRouterSetting struct {
	// Interfaces NIC設定
	Interfaces []*VPCRouterInterface `json:",omitempty"`
	// StaticNAT スタティックNAT設定
	StaticNAT *VPCRouterStaticNAT `json:",omitempty"`
	// PortForwarding ポートフォワーディング設定
	PortForwarding *VPCRouterPortForwarding `json:",omitempty"`
	// Firewall ファイアウォール設定
	Firewall *VPCRouterFirewall `json:",omitempty"`
	// DHCPServer DHCPサーバー設定
	DHCPServer *VPCRouterDHCPServer `json:",omitempty"`
	// DHCPStaticMapping DHCPスタティックマッピング設定
	DHCPStaticMapping *VPCRouterDHCPStaticMapping `json:",omitempty"`
	// L2TPIPsecServer L2TP/IPSecサーバー設定
	L2TPIPsecServer *VPCRouterL2TPIPsecServer `json:",omitempty"`
	// PPTPServer PPTPサーバー設定
	PPTPServer *VPCRouterPPTPServer `json:",omitempty"`
	// RemoteAccessUsers リモートアクセスユーザー設定
	RemoteAccessUsers *VPCRouterRemoteAccessUsers `json:",omitempty"`
	// SiteToSiteIPsecVPN サイト間VPN設定
	SiteToSiteIPsecVPN *VPCRouterSiteToSiteIPsecVPN `json:",omitempty"`
	// StaticRoutes スタティックルート設定
	StaticRoutes *VPCRouterStaticRoutes `json:",omitempty"`
	// VRID VRID
	VRID *int `json:",omitempty"`
	// SyslogHost syslog転送先ホスト
	SyslogHost string `json:",omitempty"`
}

// VPCRouterInterface NIC設定
type VPCRouterInterface struct {
	// IPAddress IPアドレスリスト
	IPAddress []string `json:",omitempty"`
	// NetworkMaskLen ネットワークマスク長
	NetworkMaskLen int `json:",omitempty"`
	// VirtualIPAddress 仮想IPアドレス
	VirtualIPAddress string `json:",omitempty"`
	// IPAliases IPエイリアス
	IPAliases []string `json:",omitempty"`
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
	// Config スタティックNAT設定
	Config []*VPCRouterStaticNATConfig `json:",omitempty"`
	// Enabled 有効/無効
	Enabled string `json:",omitempty"`
}

// VPCRouterStaticNATConfig スタティックNAT設定
type VPCRouterStaticNATConfig struct {
	// GlobalAddress グローバルIPアドレス
	GlobalAddress string `json:",omitempty"`
	// PrivateAddress プライベートIPアドレス
	PrivateAddress string `json:",omitempty"`
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
	// Config ポートフォワーディング設定
	Config []*VPCRouterPortForwardingConfig `json:",omitempty"`
	// Enabled 有効/無効
	Enabled string `json:",omitempty"`
}

// VPCRouterPortForwardingConfig ポートフォワーディング設定
type VPCRouterPortForwardingConfig struct {
	// Protocol プロトコル
	Protocol string `json:",omitempty"` // tcp/udp only
	// GlobalPort グローバル側ポート
	GlobalPort string `json:",omitempty"`
	// PrivateAddress プライベートIPアドレス
	PrivateAddress string `json:",omitempty"`
	// PrivatePort プライベート側ポート
	PrivatePort string `json:",omitempty"`
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
	// Config ファイアウォール設定
	Config []*VPCRouterFirewallSetting `json:",omitempty"`
	// Enabled 有効/無効
	Enabled string `json:",omitempty"`
}

// VPCRouterFirewallSetting ファイアウォール設定
type VPCRouterFirewallSetting struct {
	// Receive 受信ルール
	Receive []*VPCRouterFirewallRule `json:",omitempty"`
	// Send 送信ルール
	Send []*VPCRouterFirewallRule `json:",omitempty"`
}

// VPCRouterFirewallRule ファイアウォール ルール
type VPCRouterFirewallRule struct {
	// Action 許可/拒否
	Action string `json:",omitempty"`
	// Protocol プロトコル
	Protocol string `json:",omitempty"`
	// SourceNetwork 送信元ネットワーク
	SourceNetwork string `json:",omitempty"`
	// SourcePort 送信元ポート
	SourcePort string `json:",omitempty"`
	// DestinationNetwork 宛先ネットワーク
	DestinationNetwork string `json:",omitempty"`
	// DestinationPort 宛先ポート
	DestinationPort string `json:",omitempty"`
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
	// Config DHCPサーバー設定
	Config []*VPCRouterDHCPServerConfig `json:",omitempty"`
	// Enabled 有効/無効
	Enabled string `json:",omitempty"`
}

// VPCRouterDHCPServerConfig DHCPサーバー設定
type VPCRouterDHCPServerConfig struct {
	// Interface 対象NIC
	Interface string `json:",omitempty"`
	// RangeStart 割り当て範囲 開始アドレス
	RangeStart string `json:",omitempty"`
	// RangeStop 割り当て範囲 終了アドレス
	RangeStop string `json:",omitempty"`
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
	// Config DHCPスタティックマッピング設定
	Config []*VPCRouterDHCPStaticMappingConfig `json:",omitempty"`
	// Enabled 有効/無効
	Enabled string `json:",omitempty"`
}

// VPCRouterDHCPStaticMappingConfig DHCPスタティックマッピング設定
type VPCRouterDHCPStaticMappingConfig struct {
	// IPAddress 割り当てIPアドレス
	IPAddress string `json:",omitempty"`
	// MACAddress ソースMACアドレス
	MACAddress string `json:",omitempty"`
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
	// Config L2TP/IPSecサーバー設定
	Config *VPCRouterL2TPIPsecServerConfig `json:",omitempty"`
	// Enabled 有効/無効
	Enabled string `json:",omitempty"`
}

// VPCRouterL2TPIPsecServerConfig L2TP/IPSecサーバー設定
type VPCRouterL2TPIPsecServerConfig struct {
	// PreSharedSecret 事前共有シークレット
	PreSharedSecret string `json:",omitempty"`
	// RangeStart 割り当て範囲 開始IPアドレス
	RangeStart string `json:",omitempty"`
	// RangeStop 割り当て範囲 終了IPアドレス
	RangeStop string `json:",omitempty"`
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
	// Config PPTPサーバー設定
	Config *VPCRouterPPTPServerConfig `json:",omitempty"`
	// Enabled 有効/無効
	Enabled string `json:",omitempty"`
}

// VPCRouterPPTPServerConfig PPTPサーバー設定
type VPCRouterPPTPServerConfig struct {
	// RangeStart 割り当て範囲 開始IPアドレス
	RangeStart string `json:",omitempty"`
	// RangeStop 割り当て範囲 終了IPアドレス
	RangeStop string `json:",omitempty"`
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
	// Config リモートアクセスユーザー設定
	Config []*VPCRouterRemoteAccessUsersConfig `json:",omitempty"`
	// Enabled 有効/無効
	Enabled string `json:",omitempty"`
}

// VPCRouterRemoteAccessUsersConfig リモートアクセスユーザー設定
type VPCRouterRemoteAccessUsersConfig struct {
	// UserName ユーザー名
	UserName string `json:",omitempty"`
	// Password パスワード
	Password string `json:",omitempty"`
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
	// Config サイト間VPC設定
	Config []*VPCRouterSiteToSiteIPsecVPNConfig `json:",omitempty"`
	// Enabled 有効/無効
	Enabled string `json:",omitempty"`
}

// VPCRouterSiteToSiteIPsecVPNConfig サイト間VPC設定
type VPCRouterSiteToSiteIPsecVPNConfig struct {
	// LocalPrefix ローカルプレフィックス リスト
	LocalPrefix []string `json:",omitempty"`
	// Peer 対向IPアドレス
	Peer string `json:",omitempty"`
	// PreSharedSecret 事前共有シークレット
	PreSharedSecret string `json:",omitempty"`
	// RemoteID 対向ID
	RemoteID string `json:",omitempty"`
	// Routes 対向プレフィックス リスト
	Routes []string `json:",omitempty"`
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
	// Config スタティックルート設定
	Config []*VPCRouterStaticRoutesConfig `json:",omitempty"`
	// Enabled 有効/無効
	Enabled string `json:",omitempty"`
}

// VPCRouterStaticRoutesConfig スタティックルート設定
type VPCRouterStaticRoutesConfig struct {
	// Prefix プレフィックス
	Prefix string `json:",omitempty"`
	// NextHop ネクストホップ
	NextHop string `json:",omitempty"`
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
