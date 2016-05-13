package sacloud

import (
	"time"
)

// SimpleMonitor type of SimpleMonitor(CommonServiceItem)
type SimpleMonitor struct {
	*Resource
	Name         string
	Description  string                 `json:",omitempty"`
	Settings     *SimpleMonitorSettings `json:",omitempty"`
	Status       *SimpleMonitorStatus   `json:",omitempty"`
	ServiceClass string                 `json:",omitempty"`
	CreatedAt    *time.Time             `json:",omitempty"`
	ModifiedAt   *time.Time             `json:",omitempty"`
	Provider     *SimpleMonitorProvider `json:",omitempty"`
	Icon         *Icon                  `json:",omitempty"`
	Tags         []string               `json:",omitempty"`
}

// SimpleMonitorSettings type of SimpleMonitorSettings
type SimpleMonitorSettings struct {
	SimpleMonitor *SimpleMonitorSetting `json:",omitempty"`
}

// SimpleMonitorSetting type of SimpleMonitorSetting
type SimpleMonitorSetting struct {
	DelayLoop   int                       `json:",omitempty"`
	HealthCheck *SimpleMonitorHealthCheck `json:",omitempty"`
	Enabled     string                    `json:",omitempty"`
	NotifyEmail *SimpleMonitorNotify      `json:",omitempty"`
	NotifySlack *SimpleMonitorNotify      `json:",omitempty"`
}

// SimpleMonitorStatus type of CommonServiceDNSStatus
type SimpleMonitorStatus struct {
	Target string `json:",omitempty"`
}

// SimpleMonitorProvider type of CommonServiceDNSProvider
type SimpleMonitorProvider struct {
	*NumberResource
	Class        string `json:",omitempty"`
	Name         string `json:",omitempty"`
	ServiceClass string `json:",omitempty"`
}

type SimpleMonitorHealthCheck struct {
	Protocol     string `json:",omitempty"`
	Port         string `json:",omitempty"`
	Path         string `json:",omitempty"`
	Status       string `json:",omitempty"`
	QName        string `json:",omitempty"`
	ExpectedData string `json:",omitempty"`
}

type SimpleMonitorNotify struct {
	Enabled             string `json:",omitempty"`
	IncomingWebhooksURL string `json:",omitempty"`
}

// CreateNewSimpleMonitor Create new CommonServiceSimpleMonitorItem
func CreateNewSimpleMonitor(target string) *SimpleMonitor {
	return &SimpleMonitor{
		//Resource: &Resource{ID: ""},
		Name: target,
		Provider: &SimpleMonitorProvider{
			Class: "simplemon",
		},
		Status: &SimpleMonitorStatus{
			Target: target,
		},

		Settings: &SimpleMonitorSettings{
			SimpleMonitor: &SimpleMonitorSetting{
				HealthCheck: &SimpleMonitorHealthCheck{},
				Enabled:     "True",
				NotifyEmail: &SimpleMonitorNotify{
					Enabled: "False",
				},
				NotifySlack: &SimpleMonitorNotify{
					Enabled: "False",
				},
			},
		},
	}

}

func AllowSimpleMonitorHealthCheckProtocol() []string {
	return []string{"http", "https", "ping", "tcp", "dns", "ssh"}
}

func createSimpleMonitorNotifyEmail() *SimpleMonitorNotify {
	return &SimpleMonitorNotify{
		Enabled: "True",
	}
}

func createSimpleMonitorNotifySlack(incomingWebhooksURL string) *SimpleMonitorNotify {
	return &SimpleMonitorNotify{
		Enabled:             "True",
		IncomingWebhooksURL: incomingWebhooksURL,
	}

}

func (s *SimpleMonitor) SetTarget(target string) {
	s.Name = target
	s.Status.Target = target
}

func (s *SimpleMonitor) SetHealthCheckPing() {
	s.Settings.SimpleMonitor.HealthCheck = &SimpleMonitorHealthCheck{
		Protocol: "ping",
	}
}

func (s *SimpleMonitor) SetHealthCheckTCP(port string) {
	s.Settings.SimpleMonitor.HealthCheck = &SimpleMonitorHealthCheck{
		Protocol: "tcp",
		Port:     port,
	}
}

func (s *SimpleMonitor) SetHealthCheckHTTP(port string, path string, status string) {
	s.Settings.SimpleMonitor.HealthCheck = &SimpleMonitorHealthCheck{
		Protocol: "http",
		Port:     port,
		Path:     path,
		Status:   status,
	}
}

func (s *SimpleMonitor) SetHealthCheckHTTPS(port string, path string, status string) {
	s.Settings.SimpleMonitor.HealthCheck = &SimpleMonitorHealthCheck{
		Protocol: "https",
		Port:     port,
		Path:     path,
		Status:   status,
	}
}

func (s *SimpleMonitor) SetHealthCheckDNS(qname string, expectedData string) {
	s.Settings.SimpleMonitor.HealthCheck = &SimpleMonitorHealthCheck{
		Protocol:     "dns",
		QName:        qname,
		ExpectedData: expectedData,
	}
}

func (s *SimpleMonitor) SetHealthCheckSSH(port string) {
	s.Settings.SimpleMonitor.HealthCheck = &SimpleMonitorHealthCheck{
		Protocol: "ssh",
		Port:     port,
	}
}

func (s *SimpleMonitor) EnableNotifyEmail() {
	s.Settings.SimpleMonitor.NotifyEmail = createSimpleMonitorNotifyEmail()
}

func (s *SimpleMonitor) DisableNotifyEmail() {
	s.Settings.SimpleMonitor.NotifyEmail = &SimpleMonitorNotify{
		Enabled: "False",
	}
}

func (s *SimpleMonitor) EnableNofitySlack(incomingWebhooksURL string) {
	s.Settings.SimpleMonitor.NotifySlack = createSimpleMonitorNotifySlack(incomingWebhooksURL)
}

func (s *SimpleMonitor) DisableNotifySlack() {
	s.Settings.SimpleMonitor.NotifySlack = &SimpleMonitorNotify{
		Enabled: "False",
	}
}

func (s *SimpleMonitor) SetDelayLoop(loop int) {
	s.Settings.SimpleMonitor.DelayLoop = loop
}
