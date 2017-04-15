package sacloud

// SSHKey 公開鍵
type SSHKey struct {
	*Resource       // ID
	propName        // 名称
	propDescription // 説明
	propCreatedAt   // 作成日時

	PublicKey   string `json:",omitempty"` // 公開鍵
	Fingerprint string `json:",omitempty"` // フィンガープリント
}

// SSHKeyGenerated 公開鍵生成戻り値(秘密鍵のダウンロード用)
type SSHKeyGenerated struct {
	SSHKey
	PrivateKey string `json:",omitempty"` // 秘密鍵
}

func (k *SSHKey) GetPublicKey() string {
	return k.PublicKey
}

func (k *SSHKey) SetPublicKey(pKey string) {
	k.PublicKey = pKey
}

func (k *SSHKey) GetFingerpinrt() string {
	return k.Fingerprint
}

func (k *SSHKeyGenerated) GetPrivateKey() string {
	return k.PrivateKey
}
