package api

type ProductInternetAPI struct {
	*baseAPI
}

func NewProductInternetAPI(client *Client) *ProductInternetAPI {
	return &ProductInternetAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "product/internet"
			},
		},
	}
}
