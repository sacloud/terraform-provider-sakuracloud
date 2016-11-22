package sacloud

import "time"

// License ライセンス
type License struct {
	*Resource
	// Name 名称
	Name string `json:",omitempty"`
	// Description 説明
	Description string `json:",omitempty"`
	// CreatedAt 作成日時
	CreatedAt *time.Time `json:",omitempty"`
	// ModifiedAt 変更日時
	ModifiedAt *time.Time `json:",omitempty"`
	// LicenseInfo ライセンス情報
	LicenseInfo *ProductLicense `json:",omitempty"`
}
