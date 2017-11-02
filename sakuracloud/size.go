package sakuracloud

const (

	// KB 1000B
	KB int64 = 1000
	// MB 1000KB
	MB = 1000 * KB
	// GB 1000MB
	GB = 1000 * MB
	// TB 1000GB
	TB = 1000 * GB
	// PB 1000TB
	PB = 1000 * TB

	// KiB 1024B
	KiB int64 = 1024
	// MiB 1024KiB
	MiB = 1024 * KiB
	// GiB 1024MiB
	GiB = 1024 * MiB
	// TiB 1024GiB
	TiB = 1024 * GiB
	// PiB 1024TiB
	PiB = 1024 * TiB
)

func toSizeMB(sizeGB int) int {
	if sizeGB == 0 {
		return 0
	}
	sizeGB64 := int64(sizeGB)
	return int(sizeGB64 * GiB / MiB)
}

func toSizeGB(sizeMB int) int {
	if sizeMB == 0 {
		return 0
	}
	sizeMB64 := int64(sizeMB)
	return int(sizeMB64 * MiB / GiB)
}
