package sacloud

import (
	"time"
)

// Internet ルーター
type Internet struct {
	*Resource
	// Name 名称
	Name string `json:",omitempty"`
	// Description 説明
	Description string `json:",omitempty"`
	// BandWidthMbps 帯域
	BandWidthMbps int `json:",omitempty"`
	// NetworkMaskLen ネットワークマスク長
	NetworkMaskLen int `json:",omitempty"`
	// Scope スコープ
	Scope EScope `json:",omitempty"`
	// ServiceClass サービスクラス
	ServiceClass string `json:",omitempty"`
	// CreatedAt 作成日時
	CreatedAt *time.Time `json:",omitempty"`
	// Icon アイコン
	Icon *Icon `json:",omitempty"`
	// Switch スイッチ
	Switch *Switch `json:",omitempty"`
	*TagsType

	//TODO Zone
	// Zone           *Zone      `json:",omitempty"`

}
