package api

import "github.com/yamamoto-febc/libsacloud/sacloud"

type IconAPI struct {
	*baseAPI
}

func NewIconAPI(client *Client) *IconAPI {
	return &IconAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "icon"
			},
		},
	}
}

func (api *IconAPI) GetImage(id string, size string) (*sacloud.Image, error) {

	res := &sacloud.Response{}
	err := api.read(id, map[string]string{"Size": size}, res)
	if err != nil {
		return nil, err
	}
	return res.Image, nil
}
