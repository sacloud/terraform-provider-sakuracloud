package sacloud

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Database struct {
	*Appliance
	Remark   *DatabaseRemark   `json:",omitempty"`
	Settings *DatabaseSettings `json:",omitempty"`
}

type DatabaseRemark struct {
	*ApplianceRemarkBase
	DBConf  *DatabaseCommonRemarks
	Network *DatabaseRemarkNetwork
	Zone    *NumberResource
	Plan    *NumberResource
}

type DatabaseRemarkNetwork struct {
	NetworkMaskLen int    `json:",omitempty"`
	DefaultRoute   string `json:",omitempty"`
}

func (s *DatabaseRemarkNetwork) UnmarshalJSON(data []byte) error {
	targetData := strings.Replace(strings.Replace(string(data), " ", "", -1), "\n", "", -1)
	if targetData == `[]` {
		return nil
	}

	tmp := &struct {
		NetworkMaskLen int    `json:",omitempty"`
		DefaultRoute   string `json:",omitempty"`
	}{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	s.NetworkMaskLen = tmp.NetworkMaskLen
	s.DefaultRoute = tmp.DefaultRoute
	return nil
}

type DatabaseCommonRemarks struct {
	Common *DatabaseCommonRemark
}

type DatabaseCommonRemark struct {
	DatabaseName     string `json:",omitempty"`
	DatabaseRevision string `json:",omitempty"`
	DatabaseTitle    string `json:",omitempty"`
	DatabaseVersion  string `json:",omitempty"`
	ReplicaPassword  string `json:",omitempty"`
	ReplicaUser      string `json:",omitempty"`
}

type DatabaseSettings struct {
	DBConf *DatabaseSetting `json:",omitempty"`
}

type DatabaseSetting struct {
	Backup *DatabaseBackupSetting `json:",omitempty"`
	Common *DatabaseCommonSetting `json:",oitempty"`
}

type DatabaseServer struct {
	IPAddress  string `json:",omitempty"`
	Port       string `json:",omitempty"`
	Enabled    string `json:",omitempty"`
	Status     string `json:",omitempty"`
	ActiveConn string `json:",omitempty"`
}

type DatabasePlan int

var DatabasePlanMini = DatabasePlan(1)

//var DatabasePlanPremium = DatabasePlan(2)

type DatabaseBackupSetting struct {
	Rotate string `json:",omitempty"`
	Time   string `json:",omitempty"`
}

type DatabaseCommonSetting struct {
	AdminPassword string        `json:",omitempty"`
	DefaultUser   string        `json:",omitempty"`
	UserPassword  string        `json:",omitempty"`
	ServicePort   string        `json:",omitempty"`
	SourceNetwork SourceNetwork `json:",omitempty"`
}

type SourceNetwork []string

func (s *SourceNetwork) UnmarshalJSON(data []byte) error {
	// SourceNetworkが未設定の場合、APIレスポンスが""となるため回避する
	if string(data) == `""` {
		return nil
	}

	tmp := []string{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	source := SourceNetwork(tmp)
	*s = source
	return nil
}

type CreateDatabaseValue struct {
	Plan DatabasePlan

	AdminPassword string
	DefaultUser   string
	UserPassword  string
	SourceNetwork []string
	ServicePort   string
	//EnableWebUI bool

	BackupRotate int
	BackupTime   string

	SwitchID     string
	IPAddress1   string
	MaskLen      int
	DefaultRoute string

	Name        string
	Description string
	Tags        []string
	Icon        *Resource

	DatabaseName     string
	DatabaseRevision string
	DatabaseTitle    string
	DatabaseVersion  string
	//ReplicaPassword  string //set admin password
	ReplicaUser string
}

func NewCreatePostgreSQLDatabaseValue() *CreateDatabaseValue {
	return &CreateDatabaseValue{
		DatabaseName:     "postgres",
		DatabaseRevision: "9.4.7",
		DatabaseTitle:    "PostgreSQL 9.4.7",
		DatabaseVersion:  "9.4",
		ReplicaUser:      "replica",
	}
}

func CreateNewPostgreSQLDatabase(values *CreateDatabaseValue) *Database {

	db := &Database{
		Appliance: &Appliance{
			Class:       "database",
			Name:        values.Name,
			Description: values.Description,
			Tags:        values.Tags,
			Icon: &Icon{
				Resource: values.Icon,
			},
			Plan: &NumberResource{ID: json.Number(fmt.Sprintf("%d", values.Plan))},
		},
		Remark: &DatabaseRemark{
			ApplianceRemarkBase: &ApplianceRemarkBase{
				Servers: []interface{}{""},
			},
			DBConf: &DatabaseCommonRemarks{
				Common: &DatabaseCommonRemark{
					DatabaseName:     values.DatabaseName,
					DatabaseRevision: values.DatabaseRevision,
					DatabaseTitle:    values.DatabaseTitle,
					DatabaseVersion:  values.DatabaseVersion,
					ReplicaUser:      values.ReplicaUser,
					ReplicaPassword:  values.AdminPassword,
				},
			},
			Plan: &NumberResource{ID: json.Number(fmt.Sprintf("%d", values.Plan))},
		},
		Settings: &DatabaseSettings{
			DBConf: &DatabaseSetting{
				Backup: &DatabaseBackupSetting{
					Rotate: fmt.Sprintf("%d", values.BackupRotate),
					Time:   values.BackupTime,
				},
				Common: &DatabaseCommonSetting{
					AdminPassword: values.AdminPassword,
					DefaultUser:   values.DefaultUser,
					UserPassword:  values.UserPassword,
					SourceNetwork: SourceNetwork(values.SourceNetwork),
					ServicePort:   values.ServicePort,
				},
			},
		},
	}

	if values.SwitchID == "" || values.SwitchID == "shared" {
		db.Remark.Switch = &ApplianceRemarkSwitch{
			Scope: "shared",
		}
	} else {
		db.Remark.Switch = &ApplianceRemarkSwitch{
			ID: values.SwitchID,
		}
		db.Remark.Network = &DatabaseRemarkNetwork{
			NetworkMaskLen: values.MaskLen,
			DefaultRoute:   values.DefaultRoute,
		}

		db.Remark.Servers = []interface{}{
			map[string]string{"IPAddress": values.IPAddress1},
		}

	}

	return db
}

func (s *Database) AddSourceNetwork(nw string) {
	res := []string(s.Settings.DBConf.Common.SourceNetwork)
	res = append(res, nw)
	s.Settings.DBConf.Common.SourceNetwork = SourceNetwork(res)
}

func (s *Database) DeleteSourceNetwork(nw string) {
	res := []string{}
	for _, s := range s.Settings.DBConf.Common.SourceNetwork {
		if s != nw {
			res = append(res, s)
		}
	}
	s.Settings.DBConf.Common.SourceNetwork = SourceNetwork(res)
}
