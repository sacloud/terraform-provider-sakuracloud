package sacloud

import "time"

// CDROM ISOイメージ(CDROM)
type CDROM struct {
	*Resource
	// StorageClass ストレージクラス
	StorageClass string `json:",omitempty"`
	// Name 名称
	Name string `json:",omitempty"`
	// Description 説明
	Description string `json:",omitempty"`
	// SizeMB サイズ(MB単位)
	SizeMB int `json:",omitempty"`
	// Scope スコープ
	Scope string `json:",omitempty"`
	*EAvailability
	// ServiceClass サービスクラス
	ServiceClass string `json:",omitempty"`
	// CreatedAt 作成日時
	CreatedAt *time.Time `json:",omitempty"`
	// Icon アイコン
	Icon string `json:",omitempty"`
	// Storage ストレージ
	Storage *Storage `json:",omitempty"`
	*TagsType
}
