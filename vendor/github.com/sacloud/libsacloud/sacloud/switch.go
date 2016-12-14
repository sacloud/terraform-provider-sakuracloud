package sacloud

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// Switch スイッチ
type Switch struct {
	*Resource
	// Name 名称
	Name string `json:",omitempty"`
	// Description 説明
	Description string `json:",omitempty"`
	// ServerCount 接続サーバー数
	ServerCount int `json:",omitempty"`
	// ApplianceCount 接続アプライアンス数
	ApplianceCount int `json:",omitempty"`
	// Scope スコープ
	Scope EScope `json:",omitempty"`
	// Subnet サブネット
	Subnet *Subnet `json:",omitempty"`
	// UserSubnet ユーザー定義サブネット
	UserSubnet *Subnet `json:",omitempty"`
	//HybridConnection //REMARK: !!ハイブリッド接続 not support!!
	// ServerClass サービスクラス
	ServerClass string `json:",omitempty"`
	// CreatedAt 作成日時
	CreatedAt *time.Time `json:",omitempty"`
	// Icon アイコン
	Icon *Icon `json:",omitempty"`
	*TagsType
	// Subnets サブネット
	Subnets []SwitchSubnet `json:",omitempty"`
	// IPv6Nets IPv6サブネットリスト
	IPv6Nets []IPv6Net `json:",omitempty"`
	// Internet ルーター
	Internet *Internet `json:",omitempty"`
	// Bridge ブリッジ
	Bridge *struct {
		*Bridge
		Info *struct {
			// Switches 接続スイッチリスト
			Switches []struct {
				*Switch
				ID json.Number `json:",omitempty"` // HACK
			}
		}
	} `json:",omitempty"`
}

// SwitchSubnet スイッチサブネット
type SwitchSubnet struct {
	*Subnet
	// IPAddresses IPアドレス範囲
	IPAddresses struct {
		// Min IPアドレス開始
		Min string `json:",omitempty"`
		// Max IPアドレス終了
		Max string `json:",omitempty"`
	}
}

// GetDefaultIPAddressesForVPCRouter VPCルーター接続用にサブネットからIPアドレスを3つ取得
func (s *Switch) GetDefaultIPAddressesForVPCRouter() (string, string, string, error) {

	if s.Subnets == nil || len(s.Subnets) < 1 {
		return "", "", "", fmt.Errorf("switch[%d].Subnets is nil", s.ID)
	}

	baseAddress := net.ParseIP(s.Subnets[0].IPAddresses.Min).To4()
	address1 := net.IPv4(baseAddress[0], baseAddress[1], baseAddress[2], baseAddress[3]+1)
	address2 := net.IPv4(baseAddress[0], baseAddress[1], baseAddress[2], baseAddress[3]+2)

	return baseAddress.String(), address1.String(), address2.String(), nil
}

// GetIPAddressList IPアドレス範囲内の全てのIPアドレスを取得
func (s *Switch) GetIPAddressList() ([]string, error) {
	if s.Subnets == nil || len(s.Subnets) < 1 {
		return nil, fmt.Errorf("switch[%d].Subnets is nil", s.ID)
	}

	//さくらのクラウドの仕様上/24までしか割り当てできないためこのロジックでOK
	baseIP := net.ParseIP(s.Subnets[0].IPAddresses.Min).To4()
	min := baseIP[3]
	max := net.ParseIP(s.Subnets[0].IPAddresses.Max).To4()[3]

	var i byte
	ret := []string{}
	for (min + i) <= max { //境界含む
		ip := net.IPv4(baseIP[0], baseIP[1], baseIP[2], baseIP[3]+i)
		ret = append(ret, ip.String())
		i++
	}

	return ret, nil
}
