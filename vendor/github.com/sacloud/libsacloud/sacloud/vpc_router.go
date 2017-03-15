package sacloud

// VPCRouter VPCルーター
type VPCRouter struct {
	*Appliance // アプライアンス共通属性

	Remark   *VPCRouterRemark   `json:",omitempty"` // リマーク
	Settings *VPCRouterSettings `json:",omitempty"` // VPCルーター設定リスト
}

// VPCRouterRemark リマーク
type VPCRouterRemark struct {
	*ApplianceRemarkBase
	// TODO Zone
	//Zone *Resource
}

// VPCRouterSettings VPCルーター設定リスト
type VPCRouterSettings struct {
	Router *VPCRouterSetting `json:",omitempty"` // VPCルーター設定
}

// CreateNewVPCRouter VPCルーター作成
func CreateNewVPCRouter() *VPCRouter {
	return &VPCRouter{
		Appliance: &Appliance{
			Class:      "vpcrouter",
			propPlanID: propPlanID{Plan: &Resource{}},
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

// IsPremiumPlan プレミアムプランか判定
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
		propScope: propScope{Scope: "shared"},
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
