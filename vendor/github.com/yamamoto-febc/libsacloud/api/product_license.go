package api

type ProductLicenseAPI struct {
	*baseAPI
}

func NewProductLicenseAPI(client *Client) *ProductLicenseAPI {
	return &ProductLicenseAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "product/license"
			},
		},
	}
}
