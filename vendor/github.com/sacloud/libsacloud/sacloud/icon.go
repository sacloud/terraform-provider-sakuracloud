package sacloud

import "time"

// Icon アイコン
type Icon struct {
	*Resource
	// URL アイコンURL
	URL string `json:",omitempty"`
	// Name 名称
	Name string `json:",omitempty"`
	// Image 画像データBase64文字列(Sizeパラメータ指定時 or 画像アップロード時に利用)
	Image string `json:",omitempty"`
	// Scope スコープ
	Scope string `json:",omitempty"`
	*EAvailability
	// CreatedAt 作成日時
	CreatedAt *time.Time `json:",omitempty"`
	// ModifiedAt 変更日時
	ModifiedAt *time.Time `json:",omitempty"`
	*TagsType
}

// Image 画像データBASE64文字列
type Image string
