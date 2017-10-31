package ftps

import (
	"time"
)

type EntryType int

const (
	EntryTypeFile EntryType = iota
	EntryTypeFolder
	EntryTypeLink
)

type Entry struct {
	Type EntryType
	Name string
	Size uint64
	Time time.Time
}
