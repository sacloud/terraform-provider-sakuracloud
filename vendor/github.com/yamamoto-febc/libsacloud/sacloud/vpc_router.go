package sacloud

type VPCRouter struct {
	*Appliance
	Remark   *VPCRouterRemark   `json:",omitempty"`
	Settings *VPCRouterSettings `json:",omitempty"`
}

type VPCRouterRemark struct {
	*ApplianceRemarkBase
	Zone *Resource
}

type VPCRouterSettings struct {
	Router *VPCRouterSetting `json:",omitempty"`
}

type VPCRouterSetting struct {
	Interfaces []*VPCRouterInterface `json:",omitempty"`
	VRID       int                   `json:",omitempty"`
}
type VPCRouterInterface struct {
	IPAddress        []string `json:",omitempty"`
	NetworkMaskLen   int      `json:",omitempty"`
	VirtualIPAddress string   `json:",omitempty"`
}

func CreateNewVPCRouter() *VPCRouter {
	return &VPCRouter{
		Appliance: &Appliance{
			Class: "vpcrouter",
			Plan:  &NumberResource{},
		},
		Remark: &VPCRouterRemark{
			ApplianceRemarkBase: &ApplianceRemarkBase{
				Servers: []interface{}{""},
				Switch:  &ApplianceRemarkSwitch{},
			},
		},
	}
}

func (v *VPCRouter) SetStandardPlan() {
	v.Plan.SetIDByNumber(1)
	v.Remark.Switch = &ApplianceRemarkSwitch{
		Scope: "shared",
	}
	v.Settings = nil
}

func (v *VPCRouter) SetPremiumPlan(switchID string, virtualIPAddress string, ipAddress1 string, ipAddress2 string, vrid int, networkMaskLen int) {
	v.Plan.SetIDByNumber(2)
	v.setPremiumServices(switchID, virtualIPAddress, ipAddress1, ipAddress2, vrid, networkMaskLen)
}

func (v *VPCRouter) SetHighSpecPlan(switchID string, virtualIPAddress string, ipAddress1 string, ipAddress2 string, vrid int, networkMaskLen int) {
	v.Plan.SetIDByNumber(3)
	v.setPremiumServices(switchID, virtualIPAddress, ipAddress1, ipAddress2, vrid, networkMaskLen)
}

func (v *VPCRouter) setPremiumServices(switchID string, virtualIPAddress string, ipAddress1 string, ipAddress2 string, vrid int, networkMaskLen int) {
	v.Remark.Switch = &ApplianceRemarkSwitch{
		ID: switchID,
	}
	v.Remark.Servers = []interface{}{
		map[string]string{"IPAddress": ipAddress1},
		map[string]string{"IPAddress": ipAddress2},
	}

	v.Settings = &VPCRouterSettings{
		Router: &VPCRouterSetting{
			Interfaces: []*VPCRouterInterface{
				{
					IPAddress: []string{
						ipAddress1,
						ipAddress2,
					},
					VirtualIPAddress: virtualIPAddress,
					NetworkMaskLen:   networkMaskLen,
				},
			},
			VRID: vrid,
		},
	}

}
