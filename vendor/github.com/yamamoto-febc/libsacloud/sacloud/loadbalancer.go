package sacloud

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

func CreateNewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{}
}
