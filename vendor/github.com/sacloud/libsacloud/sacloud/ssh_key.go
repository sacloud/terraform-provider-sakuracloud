package sacloud

import "time"

// SSHKey 公開鍵
type SSHKey struct {
	*Resource
	// Name 名称
	Name string `json:",omitempty"`
	// Description 説明
	Description string `json:",omitempty"`
	// PublicKey 公開鍵
	PublicKey string `json:",omitempty"`
	// Fingerprint フィンガープリント
	Fingerprint string `json:",omitempty"`
	// CreatedAt 作成日時
	CreatedAt *time.Time `json:",omitempty"`
}
