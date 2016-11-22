package sacloud

// Zone ゾーン
type Zone struct {
	*Resource
	// Name 名称
	Name string `json:",omitempty"`
	// Description 説明
	Description string `json:",omitempty"`
	// IsDummy ダミーフラグ
	IsDummy bool `json:",omitempty"`
	// VNCProxy VPCプロキシ
	VNCProxy struct {
		// HostName ホスト名
		HostName string `json:",omitempty"`
		// IPAddress IPアドレス
		IPAddress string `json:",omitempty"`
	} `json:",omitempty"`
	// FTPServer FTPサーバー
	FTPServer struct {
		// HostName ホスト名
		HostName string `json:",omitempty"`
		// IPAddress IPアドレス
		IPAddress string `json:",omitempty"`
	} `json:",omitempty"`
	// Region リージョン
	Region *Region `json:",omitempty"`
}
