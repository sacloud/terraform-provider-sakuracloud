package api

import (
	"fmt"
	"github.com/yamamoto-febc/libsacloud/sacloud"
)

type ProductServerAPI struct {
	*baseAPI
}

func NewProductServerAPI(client *Client) *ProductServerAPI {
	return &ProductServerAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "product/server"
			},
		},
	}
}

func (api *ProductServerAPI) getPlanIDBySpec(core int, memGB int) (string, error) {
	//assert args
	if core <= 0 {
		return "", fmt.Errorf("Invalid Parameter: CPU Core")
	}
	if memGB <= 0 {
		return "", fmt.Errorf("Invalid Parameter: Memory Size(GB)")
	}

	return fmt.Sprintf("%d%03d", memGB, core), nil
}

// IsValidPlan return validate result
func (api *ProductServerAPI) IsValidPlan(core int, memGB int) (bool, error) {

	planID, err := api.getPlanIDBySpec(core, memGB)
	if err != nil {
		return false, err
	}
	productServer, err := api.Read(planID)

	if err != nil {
		return false, err
	}

	if productServer != nil {
		return true, nil
	}

	return false, fmt.Errorf("Server Plan[%s] Not Found", planID)

}

// GetBySpec return valid plan
func (api *ProductServerAPI) GetBySpec(core int, memGB int) (*sacloud.ProductServer, error) {
	planID, err := api.getPlanIDBySpec(core, memGB)

	productServer, err := api.Read(planID)

	if err != nil {
		return nil, err
	}

	return productServer, nil
}
