package sacloud

import (
	"encoding/json"
	"time"
)

// Bridge ブリッジ
type Bridge struct {
	*Resource
	// Name 名称
	Name string `json:",omitempty"`
	// Description 説明
	Description string `json:",omitempty"`
	// Info インフォ
	Info *struct {
		// Switches 接続スイッチリスト
		Switches []*struct {
			*Switch
			ID json.Number `json:",omitempty"` // HACK
		}
	}
	// ServiceClass サービスクラス
	ServiceClass string `json:",omitempty"`
	// CreatedAt 作成日時
	CreatedAt *time.Time `json:",omitempty"`
	// Region リージョン
	Region *Region `json:",omitempty"`
	// SwitchInZone　ゾーン内接続スイッチ
	SwitchInZone *struct {
		*Resource
		// Name 名称
		Name string `json:",omitempty"`
		// ServerCount 接続サーバー数
		ServerCount int `json:",omitempty"`
		// ApplianceCount 接続アプライアンス数
		ApplianceCount int `json:",omitempty"`
		// Scope スコープ
		Scope string `json:",omitempty"`
	}
}
