package sacloud

import "time"

// Subnet IPv4サブネット
type Subnet struct {
	*Resource
	// DefaultRoute デフォルトルート
	DefaultRoute string `json:",omitempty"`
	// CreatedAt 作成日時
	CreatedAt *time.Time `json:",omitempty"`
	// IPAddresses IPv4アドレス範囲
	IPAddresses []*IPAddress `json:",omitempty"`
	// NetworkAddress ネットワークアドレス
	NetworkAddress string `json:",omitempty"`
	// NetworkMaskLen ネットワークマスク長
	NetworkMaskLen int `json:",omitempty"`
	// ServiceClass サービスクラス
	ServiceClass string `json:",omitempty"`
	// ServiceID サービスID
	ServiceID int64 `json:",omitempty"`
	// StaticRoute スタティックルート
	StaticRoute string `json:",omitempty"`
	// Switch スイッチ
	Switch *Switch `json:",omitempty"`
	// Internet ルーター
	Internet *Internet `json:",omitempty"`
}
