package sacloud

// FTPServer FTPサーバー接続情報
type FTPServer struct {
	// HostName FTPサーバーホスト名
	HostName string `json:",omitempty"`
	// IPAddress FTPサーバー IPアドレス
	IPAddress string `json:",omitempty"`
	// User 接続ユーザー名
	User string `json:",omitempty"`
	// Password パスワード
	Password string `json:",omitempty"`
}

// FTPOpenRequest FTP接続オープンリクエスト
type FTPOpenRequest struct {
	// ChangePassword パスワード変更フラグ
	ChangePassword bool //省略不可
}
