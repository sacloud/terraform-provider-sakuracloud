// Package ostype is define OS type of SakuraCloud public archive
package ostype

//go:generate stringer -type=ArchiveOSTypes

// ArchiveOSTypes パブリックアーカイブOS種別
type ArchiveOSTypes int

const (
	// CentOS OS種別:CentOS
	CentOS ArchiveOSTypes = iota
	// Ubuntu OS種別:Ubuntu
	Ubuntu
	// Debian OS種別:Debian
	Debian
	// VyOS OS種別:VyOS
	VyOS
	// CoreOS OS種別:CoreOS
	CoreOS
	// RancherOS OS種別:RancherOS
	RancherOS
	// Kusanagi OS種別:Kusanagi(CentOS)
	Kusanagi
	// SiteGuard OS種別:SiteGuard(CentOS)
	SiteGuard
	// FreeBSD OS種別:FreeBSD
	FreeBSD
	// Windows2012 OS種別:Windows Server 2012 R2 Datacenter Edition
	Windows2012
	// Windows2012RDS OS種別:Windows Server 2012 R2 for RDS
	Windows2012RDS
	// Windows2012RDSOffice OS種別:Windows Server 2012 R2 for RDS(Office)
	Windows2012RDSOffice
	// Windows2016 OS種別:Windows Server 2016 Datacenter Edition
	Windows2016
	// Windows2016RDS OS種別:Windows Server 2016 RDS
	Windows2016RDS
	// Windows2016RDSOffice OS種別:Windows Server 2016 RDS(Office)
	Windows2016RDSOffice
	// Windows2016SQLServerWeb OS種別:Windows Server 2016 SQLServer(Web)
	Windows2016SQLServerWeb
	// Windows2016SQLServerStandard OS種別:Windows Server 2016 SQLServer(Standard)
	Windows2016SQLServerStandard
	// Custom OS種別:カスタム
	Custom
)

// IsWindows Windowsか
func (o ArchiveOSTypes) IsWindows() bool {
	switch o {
	case Windows2012, Windows2012RDS, Windows2012RDSOffice,
		Windows2016, Windows2016RDS, Windows2016RDSOffice,
		Windows2016SQLServerWeb, Windows2016SQLServerStandard:
		return true
	default:
		return false
	}
}

// IsSupportDiskEdit ディスクの修正機能をフルサポートしているか(Windowsは一部サポートのためfalseを返す)
func (o ArchiveOSTypes) IsSupportDiskEdit() bool {
	switch o {
	case CentOS, Ubuntu, Debian, VyOS, CoreOS, RancherOS, Kusanagi, SiteGuard, FreeBSD:
		return true
	default:
		return false
	}
}
