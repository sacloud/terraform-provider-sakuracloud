// Package ostype is define OS type of SakuraCloud public archive
package ostype

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
	// Kusanagi OS種別:Kusanagi(CentOS)
	Kusanagi
	// Custom OS種別:カスタム
	Custom
)
