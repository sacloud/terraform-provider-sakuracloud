package sacloud

import "time"

// IPv6Net IPv6ネットワーク(サブネット)
type IPv6Net struct {
	*Resource

	// IPv6Prefix IPv6プレフィックス
	IPv6Prefix string `json:",omitempty"`
	// IPv6PrefixLen IPv6プレフィックス長
	IPv6PrefixLen int `json:",omitempty"`
	// IPv6PrefixTail IPv6プレフィックス末尾
	IPv6PrefixTail string `json:",omitempty"`
	// IPv6Table IPv6テーブル
	IPv6Table *Resource `json:",omitempty"`
	// NamedIPv6AddrCount 名前付きIPv6アドレス数
	NamedIPv6AddrCount int `json:",omitempty"`
	// ServiceID サービスID
	ServiceID int64 `json:",omitempty"`
	// ServiceClass サービスクラス
	ServiceClass string `json:",omitempty"`
	// Scope スコープ
	Scope string `json:",omitempty"`
	// CreatedAt 作成日時
	CreatedAt *time.Time `json:",omitempty"`
	// Switch スイッチ
	Switch *Switch `json:",omitempty"`
}
