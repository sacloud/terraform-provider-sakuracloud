package sacloud

// Region リージョン
type Region struct {
	*Resource
	// Name 名称
	Name string `json:",omitempty"`
	// Description 説明
	Description string `json:",omitempty"`
	// NameServers ネームサーバー
	NameServers []string `json:",omitempty"`
}
