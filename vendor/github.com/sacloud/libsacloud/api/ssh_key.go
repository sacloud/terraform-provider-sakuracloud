package api

// SSHKeyAPI 公開鍵API
type SSHKeyAPI struct {
	*baseAPI
}

// NewSSHKeyAPI 公開鍵API作成
func NewSSHKeyAPI(client *Client) *SSHKeyAPI {
	return &SSHKeyAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "sshkey"
			},
		},
	}
}
