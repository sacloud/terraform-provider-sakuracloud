package api

type PacketFilterAPI struct {
	*baseAPI
}

func NewPacketFilterAPI(client *Client) *PacketFilterAPI {
	return &PacketFilterAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "packetfilter"
			},
		},
	}
}
