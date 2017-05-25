package sacloud

// Account さくらのクラウド アカウント
type Account struct {
	// *Resource //HACK 現状ではAPI戻り値が文字列なためパースエラーになる

	propName        // 名称
	ID       string `json:",omitempty"` // リソースID
	Class    string `json:",omitempty"` // リソースクラス
	Code     string `json:",omitempty"` // アカウントコード
}
