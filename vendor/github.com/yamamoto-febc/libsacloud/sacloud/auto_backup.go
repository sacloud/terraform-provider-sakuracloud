package sacloud

import "time"

// AutoBackup type of AutoBackup(CommonServiceItem)
type AutoBackup struct {
	*Resource
	Name         string
	Description  string              `json:",omitempty"`
	Status       *AutoBackupStatus   `json:",omitempty"`
	Provider     *AutoBackupProvider `json:",omitempty"`
	Settings     *AutoBackupSettings `json:",omitempty"`
	ServiceClass string              `json:",omitempty"`
	CreatedAt    *time.Time          `json:",omitempty"`
	ModifiedAt   *time.Time          `json:",omitempty"`
	Icon         *Icon               `json:",omitempty"`
	Tags         []string            //`json:",omitempty"`
}

// AutoBackupSettings type of AutoBackupSettings
type AutoBackupSettings struct {
	AccountID  string                `json:"AccountId,omitempty"`
	DiskID     string                `json:"DiskId,omitempty"`
	ZoneID     int                   `json:"ZoneId,omitempty"`
	ZoneName   string                `json:",omitempty"`
	Autobackup *AutoBackupRecordSets `json:",omitempty"`
}

// AutoBackupStatus type of AutoBackupStatus
type AutoBackupStatus struct {
	AccountID string `json:"AccountId,omitempty"`
	DiskID    string `json:"DiskId,omitempty"`
	ZoneID    int    `json:"ZoneId,omitempty"`
	ZoneName  string `json:",omitempty"`
}

// AutoBackupProvider type of AutoBackupProvider
type AutoBackupProvider struct {
	Class string `json:",omitempty"`
}

// CreateNewAutoBackup Create new AutoBackup(CommonServiceItem)
func CreateNewAutoBackup(backupName string, diskID string) *AutoBackup {
	return &AutoBackup{
		Resource: &Resource{ID: ""},
		Name:     backupName,
		Status: &AutoBackupStatus{
			DiskID: diskID,
		},
		Provider: &AutoBackupProvider{
			Class: "autobackup",
		},
		Settings: &AutoBackupSettings{
			Autobackup: &AutoBackupRecordSets{
				BackupSpanType: "weekdays",
			},
		},
	}
}

func AllowAutoBackupWeekdays() []string {
	return []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"}
}

func AllowAutoBackupHour() []int {
	return []int{0, 6, 12, 18}
}

// AutoBackupRecordSets type of AutoBackupRecordSets
type AutoBackupRecordSets struct {
	BackupSpanType          string
	BackupHour              int
	BackupSpanWeekdays      []string
	MaximumNumberOfArchives int
}

func (a *AutoBackup) SetBackupHour(hour int) {
	a.Settings.Autobackup.BackupHour = hour
}

func (a *AutoBackup) SetBackupSpanWeekdays(weekdays []string) {
	a.Settings.Autobackup.BackupSpanWeekdays = weekdays
}

func (a *AutoBackup) SetBackupMaximumNumberOfArchives(max int) {
	a.Settings.Autobackup.MaximumNumberOfArchives = max
}
