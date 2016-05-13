package api

type LicenseAPI struct {
	*baseAPI
}

func NewLicenseAPI(client *Client) *LicenseAPI {
	return &LicenseAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "license"
			},
		},
	}
}
