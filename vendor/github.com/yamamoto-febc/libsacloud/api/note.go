package api

import (
	"fmt"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"time"
)

const (
	sakuraAllowSudoScriptBody = `#!/bin/bash

  # @sacloud-once
  # @sacloud-desc ubuntuユーザーがsudo出来るように/etc/sudoersを編集します
  # @sacloud-desc （このスクリプトは、DebianもしくはUbuntuでのみ動作します）
  # @sacloud-require-archive distro-debian
  # @sacloud-require-archive distro-ubuntu

  export DEBIAN_FRONTEND=noninteractive
	echo "ubuntu ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers || exit 1
	sh -c 'sleep 10; shutdown -h now' &
  exit 0`

	sakuraAddIPForEth1ScriptBodyFormat = `#!/bin/bash

	# @sacloud-once
	# @sacloud-desc docker-machine-sakuracloud: setup ip address for eth1
	# @sacloud-desc （このスクリプトは、DebianもしくはUbuntuでのみ動作します）
	# @sacloud-require-archive distro-debian
	# @sacloud-require-archive distro-ubuntu

	export DEBIAN_FRONTEND=noninteractive
	echo "auto eth1" >> /etc/network/interfaces
	echo "iface eth1 inet static" >> /etc/network/interfaces
	echo "address %s" >> /etc/network/interfaces
	echo "netmask %s" >> /etc/network/interfaces
	ifdown eth1; ifup eth1
	exit 0`

	sakuraChangeDefaultGatewayScriptBody = `#!/bin/bash

	# @sacloud-once
	# @sacloud-desc docker-machine-sakuracloud: change default gateway
	# @sacloud-desc （このスクリプトは、DebianもしくはUbuntuでのみ動作します）
	# @sacloud-require-archive distro-debian
	# @sacloud-require-archive distro-ubuntu

	export DEBIAN_FRONTEND=noninteractive
	sed -i 's/gateway/#gateway/g' /etc/network/interfaces
	echo "up route add default gw %s" >> /etc/network/interfaces
	exit 0`

	sakuraDisableEth0ScriptBody = `#!/bin/bash

	# @sacloud-once
	# @sacloud-desc docker-machine-sakuracloud: disable eth0
	# @sacloud-desc （このスクリプトは、DebianもしくはUbuntuでのみ動作します）
	# @sacloud-require-archive distro-debian
	# @sacloud-require-archive distro-ubuntu

	export DEBIAN_FRONTEND=noninteractive
	sed -i 's/iface eth0 inet static/iface eth0 inet manual/g' /etc/network/interfaces
	ifdown eth0 || exit 0
	exit 0`
)

type NoteAPI struct {
	*baseAPI
}

func NewNoteAPI(client *Client) *NoteAPI {
	return &NoteAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "note"
			},
		},
	}
}

// GetAllowSudoNoteID get ubuntu customize note id
// FIXME
// workaround for [Non root ssh create sudo can't get password](https://github.com/docker/machine/issues/1569)
// [PR #1586](https://github.com/docker/machine/pull/1586)がマージされるまで暫定
// スクリプト(Note)を使ってubuntuユーザがsudo可能にする
func (api *NoteAPI) GetAllowSudoNoteID(noteNamePrefix string) (string, error) {
	noteName := fmt.Sprintf("_99_%s_%d__", noteNamePrefix, time.Now().UnixNano())
	return api.findOrCreateBy(noteName, sakuraAllowSudoScriptBody)
}

// GetAddIPCustomizeNoteID get add ip customize note id
func (api *NoteAPI) GetAddIPCustomizeNoteID(noteNamePrefix string, ip string, subnet string) (string, error) {
	noteName := fmt.Sprintf("_30_%s_%d__", noteNamePrefix, time.Now().UnixNano())
	noteBody := fmt.Sprintf(sakuraAddIPForEth1ScriptBodyFormat, ip, subnet)
	return api.findOrCreateBy(noteName, noteBody)
}

// GetChangeGatewayCustomizeNoteID get change gateway address customize note id
func (api *NoteAPI) GetChangeGatewayCustomizeNoteID(noteNamePrefix string, gateway string) (string, error) {
	noteName := fmt.Sprintf("_20_%s_%d__", noteNamePrefix, time.Now().UnixNano())
	noteBody := fmt.Sprintf(sakuraChangeDefaultGatewayScriptBody, gateway)
	return api.findOrCreateBy(noteName, noteBody)
}

// GetDisableEth0CustomizeNoteID get disable eth0 customize note id
func (api *NoteAPI) GetDisableEth0CustomizeNoteID(noteNamePrefix string) (string, error) {
	noteName := fmt.Sprintf("_10_%s_%d__", noteNamePrefix, time.Now().UnixNano())
	return api.findOrCreateBy(noteName, sakuraDisableEth0ScriptBody)
}

func (api *NoteAPI) findOrCreateBy(noteName string, noteBody string) (string, error) {

	existsNotes, err := api.WithNameLike(noteName).Find()
	if err != nil {
		return "", err
	}
	//すでに登録されている場合
	if existsNotes.Count > 0 {
		return existsNotes.Notes[0].ID, nil
	}

	note := &sacloud.Note{
		Name:    noteName,
		Content: noteBody,
	}

	res, err := api.Create(note)
	if err != nil {
		return "", err
	}

	return res.ID, nil
}
