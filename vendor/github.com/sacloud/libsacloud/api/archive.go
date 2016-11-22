package api

import (
	"fmt"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/libsacloud/sacloud/ostype"
	"strings"
	"time"
)

// ArchiveAPI アーカイブAPI
type ArchiveAPI struct {
	*baseAPI
	findFuncMapPerOSType map[ostype.ArchiveOSTypes]func() (*sacloud.Archive, error)
}

var (
	archiveLatestStableCentOSTags   = []string{"current-stable", "distro-centos"}
	archiveLatestStableUbuntuTags   = []string{"current-stable", "distro-ubuntu"}
	archiveLatestStableDebianTags   = []string{"current-stable", "distro-debian"}
	archiveLatestStableVyOSTags     = []string{"current-stable", "distro-vyos"}
	archiveLatestStableCoreOSTags   = []string{"current-stable", "distro-coreos"}
	archiveLatestStableKusanagiTags = []string{"current-stable", "pkg-kusanagi"}
	//ArchiveLatestStableSiteGuardTags = []string{"current-stable", "pkg-siteguard"} //tk1aではcurrent-stableタグが付いていないため絞り込めない
)

// NewArchiveAPI アーカイブAPI作成
func NewArchiveAPI(client *Client) *ArchiveAPI {
	api := &ArchiveAPI{
		baseAPI: &baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "archive"
			},
		},
	}

	api.findFuncMapPerOSType = map[ostype.ArchiveOSTypes]func() (*sacloud.Archive, error){
		ostype.CentOS:   api.FindLatestStableCentOS,
		ostype.Ubuntu:   api.FindLatestStableUbuntu,
		ostype.Debian:   api.FindLatestStableDebian,
		ostype.VyOS:     api.FindLatestStableVyOS,
		ostype.CoreOS:   api.FindLatestStableCoreOS,
		ostype.Kusanagi: api.FindLatestStableKusanagi,
	}

	return api
}

// OpenFTP FTP接続開始
func (api *ArchiveAPI) OpenFTP(id int64) (*sacloud.FTPServer, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/ftp", api.getResourceURL(), id)
		//body   = map[string]bool{"ChangePassword": reset}
		res = &sacloud.Response{}
	)

	result, err := api.action(method, uri, nil, res)
	if !result || err != nil {
		return nil, err
	}

	return res.FTPServer, nil
}

// CloseFTP FTP接続終了
func (api *ArchiveAPI) CloseFTP(id int64) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%d/ftp", api.getResourceURL(), id)
	)
	return api.modify(method, uri, nil)

}

// SleepWhileCopying コピー終了まで待機
func (api *ArchiveAPI) SleepWhileCopying(id int64, timeout time.Duration) error {

	current := 0 * time.Second
	interval := 5 * time.Second
	for {
		archive, err := api.Read(id)
		if err != nil {
			return err
		}

		if archive.IsAvailable() {
			return nil
		}
		time.Sleep(interval)
		current += interval

		if timeout > 0 && current > timeout {
			return fmt.Errorf("Timeout: SleepWhileCopying[disk:%d]", id)
		}
	}
}

// CanEditDisk ディスクの修正が可能か判定
func (api *ArchiveAPI) CanEditDisk(id int64) (bool, error) {

	archive, err := api.Read(id)
	if err != nil {
		return false, err
	}

	if archive == nil {
		return false, nil
	}

	// BundleInfoがあれば編集不可
	if archive.BundleInfo != nil {
		// Windows
		return false, nil
	}

	// ソースアーカイブ/ソースディスクともに持っていない場合
	if archive.SourceArchive == nil && archive.SourceDisk == nil {
		//ブランクディスクがソース
		return false, nil
	}

	for _, t := range allowDiskEditTags {
		if archive.HasTag(t) {
			// 対応OSインストール済みディスク
			return true, nil
		}
	}

	// ここまできても判定できないならソースに投げる
	if archive.SourceDisk != nil {
		return api.CanEditDisk(archive.SourceDisk.ID)
	}
	return api.client.Archive.CanEditDisk(archive.SourceArchive.ID)

}

// FindLatestStableCentOS 安定版最新のCentOSパブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableCentOS() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableCentOSTags)
}

// FindLatestStableDebian 安定版最新のDebianパブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableDebian() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableDebianTags)
}

// FindLatestStableUbuntu 安定版最新のUbuntuパブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableUbuntu() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableUbuntuTags)
}

// FindLatestStableVyOS 安定版最新のVyOSパブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableVyOS() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableVyOSTags)
}

// FindLatestStableCoreOS 安定版最新のCoreOSパブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableCoreOS() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableCoreOSTags)
}

// FindLatestStableKusanagi 安定版最新のKusanagiパブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableKusanagi() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableKusanagiTags)
}

// FindByOSType 指定のOS種別の安定版最新のパブリックアーカイブを取得
func (api *ArchiveAPI) FindByOSType(os ostype.ArchiveOSTypes) (*sacloud.Archive, error) {
	if f, ok := api.findFuncMapPerOSType[os]; ok {
		return f()
	}

	return nil, fmt.Errorf("OSType [%s] is invalid", os)
}

func (api *ArchiveAPI) findByOSTags(tags []string) (*sacloud.Archive, error) {
	res, err := api.Reset().WithTags(tags).Find()
	if err != nil {
		return nil, fmt.Errorf("Archive [%s] error : %s", strings.Join(tags, ","), err)
	}

	if len(res.Archives) == 0 {
		return nil, fmt.Errorf("Archive [%s] Not Found", strings.Join(tags, ","))
	}

	return &res.Archives[0], nil

}
