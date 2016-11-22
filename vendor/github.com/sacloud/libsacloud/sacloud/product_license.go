package sacloud

import "time"

// ProductLicense ライセンスプラン
type ProductLicense struct {
	*Resource
	// Name 名称
	Name string `json:",omitempty"`
	// ServiceClass サービスクラス
	ServiceClass string `json:",omitempty"`
	// TermsOfUse 利用規約
	TermsOfUse string `json:",omitempty"`
	// CreatedAt 作成日時
	CreatedAt *time.Time `json:",omitempty"`
	// ModifiedAt 変更日時
	ModifiedAt *time.Time `json:",omitempty"`
}
