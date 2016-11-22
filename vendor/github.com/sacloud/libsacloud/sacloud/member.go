package sacloud

// Member 会員情報
type Member struct {
	// Class クラス
	Class string `json:",omitempty"`
	// Code 会員コード
	Code string `json:",omitempty"`
	// Errors [unknown type] `json:",omitempty"`
}
