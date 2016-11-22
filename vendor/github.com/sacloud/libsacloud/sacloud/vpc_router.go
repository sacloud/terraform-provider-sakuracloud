package sacloud

// VPCRouter VPCルーター
type VPCRouter struct {
	*Appliance
	// Remark リマーク
	Remark *VPCRouterRemark `json:",omitempty"`
	// Settings VPCルーター設定リスト
	Settings *VPCRouterSettings `json:",omitempty"`
}

// VPCRouterRemark リマーク
type VPCRouterRemark struct {
	*ApplianceRemarkBase
	// TODO Zone
	//Zone *Resource
}

// VPCRouterSettings VPCルーター設定リスト
type VPCRouterSettings struct {
	// Router VPCルーター設定
	Router *VPCRouterSetting `json:",omitempty"`
}

// CreateNewVPCRouter VPCルーター作成
func CreateNewVPCRouter() *VPCRouter {
	return &VPCRouter{
		Appliance: &Appliance{
			Class:    "vpcrouter",
			Plan:     &Resource{},
			TagsType: &TagsType{},
		},
		Remark: &VPCRouterRemark{
			ApplianceRemarkBase: &ApplianceRemarkBase{
				Servers: []interface{}{""},
				Switch:  &ApplianceRemarkSwitch{},
			},
		},
		Settings: &VPCRouterSettings{
			Router: &VPCRouterSetting{},
		},
	}
}

// InitVPCRouterSetting VPCルーター設定初期化
func (v *VPCRouter) InitVPCRouterSetting() {
	settings := &VPCRouterSettings{
		Router: &VPCRouterSetting{},
	}

	if v.Settings != nil && v.Settings.Router != nil && v.Settings.Router.Interfaces != nil {
		settings.Router.Interfaces = v.Settings.Router.Interfaces
	}
	if v.Settings != nil && v.Settings.Router != nil && v.Settings.Router.VRID != nil {
		settings.Router.VRID = v.Settings.Router.VRID
	}

	v.Settings = settings
}

// IsStandardPlan スタンダードプランか判定
func (v *VPCRouter) IsStandardPlan() bool {
	return v.Plan.ID == 1
}

// IsPremiumPlan プレミアうプランか判定
func (v *VPCRouter) IsPremiumPlan() bool {
	return v.Plan.ID == 2
}

// IsHighSpecPlan ハイスペックプランか判定
func (v *VPCRouter) IsHighSpecPlan() bool {
	return v.Plan.ID == 3
}

// SetStandardPlan スタンダードプランへ設定
func (v *VPCRouter) SetStandardPlan() {
	v.Plan.SetID(1)
	v.Remark.Switch = &ApplianceRemarkSwitch{
		// Scope
		Scope: "shared",
	}
	v.Settings = nil
}

// SetPremiumPlan プレミアムプランへ設定
func (v *VPCRouter) SetPremiumPlan(switchID string, virtualIPAddress string, ipAddress1 string, ipAddress2 string, vrid int, ipAliases []string) {
	v.Plan.SetID(2)
	v.setPremiumServices(switchID, virtualIPAddress, ipAddress1, ipAddress2, vrid, ipAliases)
}

// SetHighSpecPlan ハイスペックプランへ設定
func (v *VPCRouter) SetHighSpecPlan(switchID string, virtualIPAddress string, ipAddress1 string, ipAddress2 string, vrid int, ipAliases []string) {
	v.Plan.SetID(3)
	v.setPremiumServices(switchID, virtualIPAddress, ipAddress1, ipAddress2, vrid, ipAliases)
}

func (v *VPCRouter) setPremiumServices(switchID string, virtualIPAddress string, ipAddress1 string, ipAddress2 string, vrid int, ipAliases []string) {
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
					IPAliases:        ipAliases,
				},
			},
			VRID: &vrid,
		},
	}

}
