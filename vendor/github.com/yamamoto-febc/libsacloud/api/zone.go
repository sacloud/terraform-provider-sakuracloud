package api

type ZoneAPI struct {
	*baseAPI
}

func NewZoneAPI(client *Client) *ZoneAPI {
	return &ZoneAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "zone"
			},
		},
	}
}
