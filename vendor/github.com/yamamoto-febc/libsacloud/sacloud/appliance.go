package sacloud

import "time"

type Appliance struct {
	*Resource
	Class       string `json:",omitempty"`
	Name        string `json:",omitempty"`
	Description string `json:",omitempty"`
	Plan        *NumberResource
	//Settings
	SettingHash string `json:",omitempty"`
	//Remark      *ApplianceRemark `json:",omitempty"`
	*EAvailability
	Instance *struct {
		Status          string     `json:",omitempty"`
		StatusChangedAt *time.Time `json:",omitempty"`
	} `json:",omitempty"`
	ServiceClass string      `json:",omitempty"`
	CreatedAt    *time.Time  `json:",omitempty"`
	Icon         *Icon       `json:",omitempty"`
	Switch       *Switch     `json:",omitempty"`
	Interfaces   []Interface `json:",omitempty"`
	Tags         []string    `json:",omitempty"`
}

//HACK Appliance:Zone.IDがRoute/LoadBalancerの場合でデータ型が異なるため
//それぞれのstruct定義でZoneだけ上書きした構造体を定義して使う

type ApplianceRemarkBase struct {
	Servers []interface{}
	Switch  *ApplianceRemarkSwitch `json:",omitempty"`
	//Zone *NumberResource `json:",omitempty"`
	VRRP    *ApplianceRemarkVRRP    `json:",omitempty"`
	Network *ApplianceRemarkNetwork `json:",omitempty"`
	//Plan    *NumberResource
}

type ApplianceRemarkSwitch struct {
	ID    string `json:",omitempty"`
	Scope string `json:",omitempty"`
}

type ApplianceRemarkVRRP struct {
	VRID int `json:",omitempty"`
}

type ApplianceRemarkNetwork struct {
	NetworkMaskLen int    `json:",omitempty"`
	DefaultRoute   string `json:",omitempty"`
}
