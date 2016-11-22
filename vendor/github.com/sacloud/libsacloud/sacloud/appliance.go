package sacloud

import "time"

// Appliance アプライアンス基底クラス
type Appliance struct {
	*Resource
	// Class リソースクラス
	Class string `json:",omitempty"`
	// Name 名称
	Name string `json:",omitempty"`
	// Description 説明
	Description string `json:",omitempty"`
	// Plan プラン
	Plan *Resource

	//Settings

	// SettingHash 設定ハッシュ値
	SettingHash string `json:",omitempty"`

	//Remark      *ApplianceRemark `json:",omitempty"`

	*EAvailability
	// Instance インスタンス
	Instance *EServerInstanceStatus `json:",omitempty"`
	// ServiceClass サービスクラス
	ServiceClass string `json:",omitempty"`
	// CreatedAt 作成日時
	CreatedAt *time.Time `json:",omitempty"`
	// Icon アイコン
	Icon *Icon `json:",omitempty"`
	// Switch スイッチ
	Switch *Switch `json:",omitempty"`
	// Interfaces インターフェース
	Interfaces []Interface `json:",omitempty"`

	*TagsType
}

//HACK Appliance:Zone.IDがRoute/LoadBalancerの場合でデータ型が異なるため
//それぞれのstruct定義でZoneだけ上書きした構造体を定義して使う

// ApplianceRemarkBase アプライアンス Remark 基底クラス
type ApplianceRemarkBase struct {
	// Servers 配下のサーバー群
	Servers []interface{}
	// Switch 接続先スイッチ
	Switch *ApplianceRemarkSwitch `json:",omitempty"`
	//Zone *Resource `json:",omitempty"`

	// VRRP VRRP
	VRRP *ApplianceRemarkVRRP `json:",omitempty"`
	// Network ネットワーク
	Network *ApplianceRemarkNetwork `json:",omitempty"`

	//Plan    *Resource
}

//type ApplianceServer struct {
//	IPAddress string `json:",omitempty"`
//}

// ApplianceRemarkSwitch スイッチ
type ApplianceRemarkSwitch struct {
	// ID リソースID
	ID string `json:",omitempty"`
	// Scope スコープ
	Scope string `json:",omitempty"`
}

// ApplianceRemarkVRRP VRRP
type ApplianceRemarkVRRP struct {
	// VRID VRID
	VRID int `json:",omitempty"`
}

// ApplianceRemarkNetwork ネットワーク
type ApplianceRemarkNetwork struct {
	// NetworkMaskLen ネットワークマスク長
	NetworkMaskLen int `json:",omitempty"`
	// DefaultRoute デフォルトルート
	DefaultRoute string `json:",omitempty"`
}
