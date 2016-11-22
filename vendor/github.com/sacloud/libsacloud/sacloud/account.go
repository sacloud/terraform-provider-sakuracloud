package sacloud

// Account さくらのクラウド アカウント
type Account struct {
	// *Resource //HACK 現状ではAPI戻り値が文字列なためパースエラーになる

	// ID リソースID
	ID string `json:",omitempty"`
	// Class リソースクラス
	Class string `json:",omitempty"`
	// Code アカウントコード
	Code string `json:",omitempty"`
	// Name アカウント名称
	Name string `json:",omitempty"`
}
