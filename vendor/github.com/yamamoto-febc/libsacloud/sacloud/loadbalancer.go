package sacloud

import (
	"encoding/json"
	"fmt"
)

type LoadBalancer struct {
	*Appliance
	Remark   *LoadBalancerRemark   `json:",omitempty"`
	Settings *LoadBalancerSettings `json:",omitempty"`
}

type LoadBalancerRemark struct {
	*ApplianceRemarkBase
	Zone *NumberResource
}

type LoadBalancerSettings struct {
	LoadBalancer []*LoadBalancerSetting `json:",omitempty"`
}

type LoadBalancerSetting struct {
	VirtualIPAddress string                `json:",omitempty"`
	Port             string                `json:",omitempty"`
	DelayLoop        string                `json:",omitempty"`
	SorryServer      string                `json:",omitempty"`
	Servers          []*LoadBalancerServer `json:",omitempty"`
}

type LoadBalancerServer struct {
	IPAddress   string                   `json:",omitempty"`
	Port        string                   `json:",omitempty"`
	HealthCheck *LoadBalancerHealthCheck `json:",omitempty"`
	Enabled     string                   `json:",omitempty"`
	Status      string                   `json:",omitempty"`
	ActiveConn  string                   `json:",omitempty"`
}

type LoadBalancerHealthCheck struct {
	Protocol string `json:",omitempty"`
	Path     string `json:",omitempty"`
	Status   string `json:",omitempty"`
}

type LoadBalancerPlan int

var LoadBalancerPlanStandard = LoadBalancerPlan(1)
var LoadBalancerPlanPremium = LoadBalancerPlan(2)

type CreateLoadBalancerValue struct {
	SwitchID     string
	VRID         int
	Plan         LoadBalancerPlan
	IPAddress1   string
	MaskLen      int
	DefaultRoute string
	Name         string
	Description  string
	Tags         []string
	Icon         *Resource
}

type CreateDoubleLoadBalancerValue struct {
	*CreateLoadBalancerValue
	IPAddress2 string
}

func AllowLoadBalancerHealthCheckProtocol() []string {
	return []string{"http", "https", "ping", "tcp"}
}

func CreateNewLoadBalancerSingle(values *CreateLoadBalancerValue, settings []*LoadBalancerSetting) (*LoadBalancer, error) {

	lb := &LoadBalancer{
		Appliance: &Appliance{
			Class:       "loadbalancer",
			Name:        values.Name,
			Description: values.Description,
			Tags:        values.Tags,
			Plan:        &NumberResource{ID: json.Number(fmt.Sprintf("%d", values.Plan))},
			Icon: &Icon{
				Resource: values.Icon,
			},
		},
		Remark: &LoadBalancerRemark{
			ApplianceRemarkBase: &ApplianceRemarkBase{
				Switch: &ApplianceRemarkSwitch{
					ID: values.SwitchID,
				},
				VRRP: &ApplianceRemarkVRRP{
					VRID: values.VRID,
				},
				Network: &ApplianceRemarkNetwork{
					NetworkMaskLen: values.MaskLen,
					DefaultRoute:   values.DefaultRoute,
				},
				Servers: []interface{}{
					map[string]string{"IPAddress": values.IPAddress1},
				},
			},
		},
	}

	for _, s := range settings {
		lb.AddLoadBalancerSetting(s)
	}

	return lb, nil
}
func CreateNewLoadBalancerDouble(values *CreateDoubleLoadBalancerValue, settings []*LoadBalancerSetting) (*LoadBalancer, error) {
	lb, err := CreateNewLoadBalancerSingle(values.CreateLoadBalancerValue, settings)
	if err != nil {
		return nil, err
	}
	lb.Remark.Servers = append(lb.Remark.Servers, map[string]string{"IPAddress": values.IPAddress2})
	return lb, nil
}

func (l *LoadBalancer) AddLoadBalancerSetting(setting *LoadBalancerSetting) {
	if l.Settings == nil {
		l.Settings = &LoadBalancerSettings{}
	}
	if l.Settings.LoadBalancer == nil {
		l.Settings.LoadBalancer = []*LoadBalancerSetting{}
	}
	l.Settings.LoadBalancer = append(l.Settings.LoadBalancer, setting)
}

func (l *LoadBalancer) DeleteLoadBalancerSetting(vip string, port string) {
	res := []*LoadBalancerSetting{}
	for _, l := range l.Settings.LoadBalancer {
		if l.VirtualIPAddress != vip || l.Port != port {
			res = append(res, l)
		}
	}

	l.Settings.LoadBalancer = res
}

func (s *LoadBalancerSetting) AddServer(server *LoadBalancerServer) {
	if s.Servers == nil {
		s.Servers = []*LoadBalancerServer{}
	}
	s.Servers = append(s.Servers, server)
}

func (s *LoadBalancerSetting) DeleteServer(ip string, port string) {
	res := []*LoadBalancerServer{}
	for _, server := range s.Servers {
		if server.IPAddress != ip || server.Port != port {
			res = append(res, server)
		}
	}

	s.Servers = res

}
