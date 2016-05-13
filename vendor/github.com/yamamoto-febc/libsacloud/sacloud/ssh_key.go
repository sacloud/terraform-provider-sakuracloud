package sacloud

import "time"

// SSHKey type of sshkey
type SSHKey struct {
	*Resource
	Name        string     `json:",omitempty"`
	Description string     `json:",omitempty"`
	PublicKey   string     `json:",omitempty"`
	Fingerprint string     `json:",omitempty"`
	CreatedAt   *time.Time `json:",omitempty"`
}
