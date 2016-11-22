package sacloud

import "time"

// Note スタートアップスクリプト
type Note struct {
	*Resource
	// Name 名称
	Name string
	// Class クラス
	Class string `json:",omitempty"`
	// Scope スコープ
	Scope string `json:",omitempty"`
	// Content スクリプト本体
	Content string `json:",omitempty"`
	// Description 説明
	Description string `json:",omitempty"`
	*EAvailability
	// CreatedAt 作成日時
	CreatedAt *time.Time `json:",omitempty"`
	// ModifiedAt 変更日時
	ModifiedAt *time.Time `json:",omitempty"`
	// Icon アイコン
	Icon *Icon `json:",omitempty"`
	*TagsType
	//TODO Remarkオブジェクトのパース
}
