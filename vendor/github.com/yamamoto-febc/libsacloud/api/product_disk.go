package api

type ProductDiskAPI struct {
	*baseAPI
}

func NewProductDiskAPI(client *Client) *ProductDiskAPI {
	return &ProductDiskAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "product/disk"
			},
		},
	}
}
