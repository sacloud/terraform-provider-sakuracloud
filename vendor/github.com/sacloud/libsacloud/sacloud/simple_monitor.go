package sacloud

import (
	"time"
)

// SimpleMonitor シンプル監視
type SimpleMonitor struct {
	*Resource
	// Name 名称
	Name string
	// Description 説明
	Description string `json:",omitempty"`
	// Settings 設定
	Settings *SimpleMonitorSettings `json:",omitempty"`
	// Status ステータス
	Status *SimpleMonitorStatus `json:",omitempty"`
	// ServiceClass サービスクラス
	ServiceClass string `json:",omitempty"`
	// CreatedAt 作成日時
	CreatedAt *time.Time `json:",omitempty"`
	// ModifiedAt 変更日時
	ModifiedAt *time.Time `json:",omitempty"`
	// Provider プロバイダ
	Provider *SimpleMonitorProvider `json:",omitempty"`
	// Icon アイコン
	Icon *Icon `json:",omitempty"`
	*TagsType
}

// SimpleMonitorSettings シンプル監視設定 リスト
type SimpleMonitorSettings struct {
	// SimpleMonitor シンプル監視設定値
	SimpleMonitor *SimpleMonitorSetting `json:",omitempty"`
}

// SimpleMonitorSetting シンプル監視設定
type SimpleMonitorSetting struct {
	// DelayLoop 監視間隔
	DelayLoop int `json:",omitempty"`
	// HealthCheck ヘルスチェック
	HealthCheck *SimpleMonitorHealthCheck `json:",omitempty"`
	// Enabled 有効/無効
	Enabled string `json:",omitempty"`
	// NotifyEmail Email通知
	NotifyEmail *SimpleMonitorNotify `json:",omitempty"`
	// NotifySlack Slack通知
	NotifySlack *SimpleMonitorNotify `json:",omitempty"`
}

// SimpleMonitorStatus シンプル監視ステータス
type SimpleMonitorStatus struct {
	// Target 対象(IP or FQDN)
	Target string `json:",omitempty"`
}

// SimpleMonitorProvider プロバイダ
type SimpleMonitorProvider struct {
	*Resource
	// Class クラス
	Class string `json:",omitempty"`
	// Name 名称
	Name string `json:",omitempty"`
	// ServiceClass サービスクラス
	ServiceClass string `json:",omitempty"`
}

// SimpleMonitorHealthCheck ヘルスチェック
type SimpleMonitorHealthCheck struct {
	// Protocol プロトコル
	Protocol string `json:",omitempty"`
	// Port ポート
	Port string `json:",omitempty"`
	// Path HTTP/HTTPS監視の場合のリクエストパス
	Path string `json:",omitempty"`
	// Status HTTP/HTTPS監視の場合の期待ステータスコード
	Status string `json:",omitempty"`
	// Host 対象ホスト(IP or FQDN)
	Host string `json:",omitempty"`
	// QName DNS監視の場合の問い合わせFQDN
	QName string `json:",omitempty"`
	// ExpectedData 期待値
	ExpectedData string `json:",omitempty"`
	// Community SNMP監視の場合のコミュニティ名
	Community string `json:",omitempty"`
	// SNMPVersion SNMP監視 SNMPバージョン
	SNMPVersion string `json:",omitempty"`
	// OID SNMP監視 OID
	OID string `json:",omitempty"`
}

// SimpleMonitorNotify シンプル監視通知
type SimpleMonitorNotify struct {
	// Enabled 有効/無効
	Enabled string `json:",omitempty"`
	// HTML メール通知の場合のHTMLメール有効フラグ
	HTML string `json:",omitempty"`
	// IncomingWebhooksURL Slack通知の場合のWebhook URL
	IncomingWebhooksURL string `json:",omitempty"`
}

// CreateNewSimpleMonitor シンプル監視作成
func CreateNewSimpleMonitor(target string) *SimpleMonitor {
	return &SimpleMonitor{
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
		TagsType: &TagsType{},
	}

}

// AllowSimpleMonitorHealthCheckProtocol シンプル監視対応プロトコルリスト
func AllowSimpleMonitorHealthCheckProtocol() []string {
	return []string{"http", "https", "ping", "tcp", "dns", "ssh", "smtp", "pop3", "snmp"}
}

func createSimpleMonitorNotifyEmail(withHTML bool) *SimpleMonitorNotify {
	n := &SimpleMonitorNotify{
		Enabled: "True",
		HTML:    "False",
	}

	if withHTML {
		n.HTML = "True"
	}

	return n
}

func createSimpleMonitorNotifySlack(incomingWebhooksURL string) *SimpleMonitorNotify {
	return &SimpleMonitorNotify{
		Enabled:             "True",
		IncomingWebhooksURL: incomingWebhooksURL,
	}

}

// SetTarget 対象ホスト(IP or FQDN)の設定
func (s *SimpleMonitor) SetTarget(target string) {
	s.Name = target
	s.Status.Target = target
}

// SetHealthCheckPing pingでのヘルスチェック設定
func (s *SimpleMonitor) SetHealthCheckPing() {
	s.Settings.SimpleMonitor.HealthCheck = &SimpleMonitorHealthCheck{
		Protocol: "ping",
	}
}

// SetHealthCheckTCP TCPでのヘルスチェック設定
func (s *SimpleMonitor) SetHealthCheckTCP(port string) {
	s.Settings.SimpleMonitor.HealthCheck = &SimpleMonitorHealthCheck{
		Protocol: "tcp",
		Port:     port,
	}
}

// SetHealthCheckHTTP HTTPでのヘルスチェック設定
func (s *SimpleMonitor) SetHealthCheckHTTP(port string, path string, status string, host string) {
	s.Settings.SimpleMonitor.HealthCheck = &SimpleMonitorHealthCheck{
		Protocol: "http",
		Port:     port,
		Path:     path,
		Status:   status,
		Host:     host,
	}
}

// SetHealthCheckHTTPS HTTPSでのヘルスチェック設定
func (s *SimpleMonitor) SetHealthCheckHTTPS(port string, path string, status string, host string) {
	s.Settings.SimpleMonitor.HealthCheck = &SimpleMonitorHealthCheck{
		Protocol: "https",
		Port:     port,
		Path:     path,
		Status:   status,
		Host:     host,
	}
}

// SetHealthCheckDNS DNSクエリでのヘルスチェック設定
func (s *SimpleMonitor) SetHealthCheckDNS(qname string, expectedData string) {
	s.Settings.SimpleMonitor.HealthCheck = &SimpleMonitorHealthCheck{
		Protocol:     "dns",
		QName:        qname,
		ExpectedData: expectedData,
	}
}

// SetHealthCheckSSH SSHヘルスチェック設定
func (s *SimpleMonitor) SetHealthCheckSSH(port string) {
	s.Settings.SimpleMonitor.HealthCheck = &SimpleMonitorHealthCheck{
		Protocol: "ssh",
		Port:     port,
	}
}

// SetHealthCheckSMTP SMTPヘルスチェック設定
func (s *SimpleMonitor) SetHealthCheckSMTP(port string) {
	s.Settings.SimpleMonitor.HealthCheck = &SimpleMonitorHealthCheck{
		Protocol: "smtp",
		Port:     port,
	}
}

// SetHealthCheckPOP3 POP3ヘルスチェック設定
func (s *SimpleMonitor) SetHealthCheckPOP3(port string) {
	s.Settings.SimpleMonitor.HealthCheck = &SimpleMonitorHealthCheck{
		Protocol: "pop3",
		Port:     port,
	}
}

// SetHealthCheckSNMP SNMPヘルスチェック設定
func (s *SimpleMonitor) SetHealthCheckSNMP(community string, version string, oid string, expectedData string) {
	s.Settings.SimpleMonitor.HealthCheck = &SimpleMonitorHealthCheck{
		Protocol:     "snmp",
		Community:    community,
		SNMPVersion:  version,
		OID:          oid,
		ExpectedData: expectedData,
	}
}

// EnableNotifyEmail Email通知の有効か
func (s *SimpleMonitor) EnableNotifyEmail(withHTML bool) {
	s.Settings.SimpleMonitor.NotifyEmail = createSimpleMonitorNotifyEmail(withHTML)
}

// DisableNotifyEmail Email通知の無効化
func (s *SimpleMonitor) DisableNotifyEmail() {
	s.Settings.SimpleMonitor.NotifyEmail = &SimpleMonitorNotify{
		Enabled: "False",
	}
}

// EnableNofitySlack Slack通知の有効化
func (s *SimpleMonitor) EnableNofitySlack(incomingWebhooksURL string) {
	s.Settings.SimpleMonitor.NotifySlack = createSimpleMonitorNotifySlack(incomingWebhooksURL)
}

// DisableNotifySlack Slack通知の無効化
func (s *SimpleMonitor) DisableNotifySlack() {
	s.Settings.SimpleMonitor.NotifySlack = &SimpleMonitorNotify{
		Enabled: "False",
	}
}

// SetDelayLoop 監視間隔の設定
func (s *SimpleMonitor) SetDelayLoop(loop int) {
	s.Settings.SimpleMonitor.DelayLoop = loop
}
