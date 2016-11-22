package sacloud

import (
	"math"
	"time"
)

// MonitorValue アクティビティモニター
type MonitorValue struct {
	// CPUTime CPU時間
	CPUTime *float64 `json:"CPU-TIME,omitempty"`
	// Write ディスク書き込み
	Write *float64 `json:",omitempty"`
	// Read ディスク読み取り
	Read *float64 `json:",omitempty"`
	// Receive パケット受信
	Receive *float64 `json:",omitempty"`
	// Send パケット送信
	Send *float64 `json:",omitempty"`
	// In パケット受信
	In *float64 `json:",omitempty"`
	// Out パケット送信
	Out *float64 `json:",omitempty"`
	// TotalMemorySize 総メモリサイズ
	TotalMemorySize *float64 `json:"Total-Memory-Size,omitempty"`
	// UsedMemorySize 使用済みメモリサイズ
	UsedMemorySize *float64 `json:"Used-Memory-Size,omitempty"`
	// TotalDisk1Size 総ディスクサイズ
	TotalDisk1Size *float64 `json:"Total-Disk1-Size,omitempty"`
	// UsedDisk1Size 使用済みディスクサイズ
	UsedDisk1Size *float64 `json:"Used-Disk1-Size,omitempty"`
	// TotalDisk2Size 総ディスクサイズ
	TotalDisk2Size *float64 `json:"Total-Disk2-Size,omitempty"`
	// UsedDisk2Size 使用済みディスクサイズ
	UsedDisk2Size *float64 `json:"Used-Disk2-Size,omitempty"`
}

// ResourceMonitorRequest アクティビティモニター取得リクエスト
type ResourceMonitorRequest struct {
	// Start 取得開始時間
	Start *time.Time `json:",omitempty"`
	// End 取得終了時間
	End *time.Time `json:",omitempty"`
}

// NewResourceMonitorRequest アクティビティモニター取得リクエスト作成
func NewResourceMonitorRequest(start *time.Time, end *time.Time) *ResourceMonitorRequest {
	res := &ResourceMonitorRequest{}
	if start != nil {
		t := start.Truncate(time.Second)
		res.Start = &t
	}
	if end != nil {
		t := end.Truncate(time.Second)
		res.End = &t
	}
	return res
}

// ResourceMonitorResponse アクティビティモニターレスポンス
type ResourceMonitorResponse struct {
	// Data メトリクス
	Data *MonitorValues `json:",omitempty"`
}

// MonitorSummaryData メトリクスサマリー
type MonitorSummaryData struct {
	// Max 最大値
	Max float64
	// Min 最小値
	Min float64
	// Avg 平均値
	Avg float64
	// Count データ個数
	Count float64
}

// MonitorSummary アクティビティーモニター サマリー
type MonitorSummary struct {
	// CPU CPU時間サマリー
	CPU *MonitorSummaryData
	// Disk ディスク利用サマリー
	Disk *struct {
		// Write ディスク書き込みサマリー
		Write *MonitorSummaryData
		// Read ディスク読み取りサマリー
		Read *MonitorSummaryData
	}
	// Interface NIC送受信サマリー
	Interface *struct {
		// Receive 受信パケットサマリー
		Receive *MonitorSummaryData
		// Send 送信パケットサマリー
		Send *MonitorSummaryData
	}
}

// MonitorValues メトリクス リスト
type MonitorValues map[string]*MonitorValue

// FlatMonitorValue フラット化したメトリクス
type FlatMonitorValue struct {
	// Time 対象時刻
	Time time.Time
	// Value 値
	Value float64
}

// Calc サマリー計算
func (m *MonitorValues) Calc() *MonitorSummary {

	res := &MonitorSummary{}
	res.CPU = m.calcBy(func(v *MonitorValue) *float64 { return v.CPUTime })
	res.Disk = &struct {
		Write *MonitorSummaryData
		Read  *MonitorSummaryData
	}{
		Write: m.calcBy(func(v *MonitorValue) *float64 { return v.Write }),
		Read:  m.calcBy(func(v *MonitorValue) *float64 { return v.Read }),
	}
	res.Interface = &struct {
		Receive *MonitorSummaryData
		Send    *MonitorSummaryData
	}{
		Receive: m.calcBy(func(v *MonitorValue) *float64 { return v.Receive }),
		Send:    m.calcBy(func(v *MonitorValue) *float64 { return v.Send }),
	}

	return res
}

func (m *MonitorValues) calcBy(f func(m *MonitorValue) *float64) *MonitorSummaryData {
	res := &MonitorSummaryData{}
	var sum float64
	for _, data := range map[string]*MonitorValue(*m) {
		value := f(data)
		if value != nil {
			res.Count++
			res.Min = math.Min(res.Min, *value)
			res.Max = math.Max(res.Max, *value)
			sum += *value
		}
	}
	if sum > 0 && res.Count > 0 {
		res.Avg = sum / res.Count
	}

	return res
}

// FlattenCPUTimeValue フラット化 CPU時間
func (m *MonitorValues) FlattenCPUTimeValue() ([]FlatMonitorValue, error) {
	return m.flattenValue(func(v *MonitorValue) *float64 { return v.CPUTime })
}

// FlattenDiskWriteValue フラット化 ディスク書き込み
func (m *MonitorValues) FlattenDiskWriteValue() ([]FlatMonitorValue, error) {
	return m.flattenValue(func(v *MonitorValue) *float64 { return v.Write })
}

// FlattenDiskReadValue フラット化 ディスク読み取り
func (m *MonitorValues) FlattenDiskReadValue() ([]FlatMonitorValue, error) {
	return m.flattenValue(func(v *MonitorValue) *float64 { return v.Read })
}

// FlattenPacketSendValue フラット化 パケット送信
func (m *MonitorValues) FlattenPacketSendValue() ([]FlatMonitorValue, error) {
	return m.flattenValue(func(v *MonitorValue) *float64 { return v.Send })
}

// FlattenPacketReceiveValue フラット化 パケット受信
func (m *MonitorValues) FlattenPacketReceiveValue() ([]FlatMonitorValue, error) {
	return m.flattenValue(func(v *MonitorValue) *float64 { return v.Receive })
}

// FlattenInternetInValue フラット化 パケット受信
func (m *MonitorValues) FlattenInternetInValue() ([]FlatMonitorValue, error) {
	return m.flattenValue(func(v *MonitorValue) *float64 { return v.In })
}

// FlattenInternetOutValue フラット化 パケット送信
func (m *MonitorValues) FlattenInternetOutValue() ([]FlatMonitorValue, error) {
	return m.flattenValue(func(v *MonitorValue) *float64 { return v.Out })
}

// FlattenTotalMemorySizeValue フラット化 総メモリサイズ
func (m *MonitorValues) FlattenTotalMemorySizeValue() ([]FlatMonitorValue, error) {
	return m.flattenValue(func(v *MonitorValue) *float64 { return v.TotalMemorySize })
}

// FlattenUsedMemorySizeValue フラット化 使用済みメモリサイズ
func (m *MonitorValues) FlattenUsedMemorySizeValue() ([]FlatMonitorValue, error) {
	return m.flattenValue(func(v *MonitorValue) *float64 { return v.UsedMemorySize })
}

// FlattenTotalDisk1SizeValue フラット化 総ディスクサイズ
func (m *MonitorValues) FlattenTotalDisk1SizeValue() ([]FlatMonitorValue, error) {
	return m.flattenValue(func(v *MonitorValue) *float64 { return v.TotalDisk1Size })
}

// FlattenUsedDisk1SizeValue フラット化 使用済みディスクサイズ
func (m *MonitorValues) FlattenUsedDisk1SizeValue() ([]FlatMonitorValue, error) {
	return m.flattenValue(func(v *MonitorValue) *float64 { return v.UsedDisk1Size })
}

// FlattenTotalDisk2SizeValue フラット化 総ディスクサイズ
func (m *MonitorValues) FlattenTotalDisk2SizeValue() ([]FlatMonitorValue, error) {
	return m.flattenValue(func(v *MonitorValue) *float64 { return v.TotalDisk2Size })
}

// FlattenUsedDisk2SizeValue フラット化 使用済みディスクサイズ
func (m *MonitorValues) FlattenUsedDisk2SizeValue() ([]FlatMonitorValue, error) {
	return m.flattenValue(func(v *MonitorValue) *float64 { return v.UsedDisk2Size })
}

func (m *MonitorValues) flattenValue(f func(*MonitorValue) *float64) ([]FlatMonitorValue, error) {
	var res []FlatMonitorValue

	for k, v := range map[string]*MonitorValue(*m) {
		if f(v) == nil {
			continue
		}
		time, err := time.Parse(time.RFC3339, k) // RFC3339 ≒ ISO8601
		if err != nil {
			return res, err
		}
		res = append(res, FlatMonitorValue{
			// Time
			Time: time,
			// Value
			Value: *f(v),
		})
	}
	return res, nil
}

// HasValue 取得したアクティビティーモニターに有効値が含まれるか判定
func (m *MonitorValue) HasValue() bool {
	values := []*float64{
		m.CPUTime,
		m.Read, m.Receive,
		m.Send, m.Write,
		m.In, m.Out,
		m.TotalMemorySize, m.UsedMemorySize,
		m.TotalDisk1Size, m.UsedDisk1Size,
		m.TotalDisk2Size, m.UsedDisk2Size,
	}
	for _, v := range values {
		if v != nil {
			return true
		}
	}
	return false
}
