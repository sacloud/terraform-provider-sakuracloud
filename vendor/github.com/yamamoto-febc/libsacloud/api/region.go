package api

type RegionAPI struct {
	*baseAPI
}

func NewRegionAPI(client *Client) *RegionAPI {
	return &RegionAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "region"
			},
		},
	}
}
