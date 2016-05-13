package api

type PublicPriceAPI struct {
	*baseAPI
}

func NewPublicPriceAPI(client *Client) *PublicPriceAPI {
	return &PublicPriceAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "public/price"
			},
		},
	}
}
