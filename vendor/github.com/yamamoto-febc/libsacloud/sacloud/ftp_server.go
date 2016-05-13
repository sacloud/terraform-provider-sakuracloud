package sacloud

type FTPServer struct {
	HostName  string `json:",omitempty"`
	IPAddress string `json:",omitempty"`
	User      string `json:",omitempty"`
	Password  string `json:",omitempty"`
}

type FTPOpenRequest struct {
	ChangePassword bool //省略不可
}
