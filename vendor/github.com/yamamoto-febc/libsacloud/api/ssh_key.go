package api

type SSHKeyAPI struct {
	*baseAPI
}

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
