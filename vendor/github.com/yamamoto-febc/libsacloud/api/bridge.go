package api

type BridgeAPI struct {
	*baseAPI
}

func NewBridgeAPI(client *Client) *BridgeAPI {
	return &BridgeAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "bridge"
			},
		},
	}
}
