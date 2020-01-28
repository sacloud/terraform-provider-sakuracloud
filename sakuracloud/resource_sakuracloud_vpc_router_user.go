// Copyright 2016-2020 terraform-provider-sakuracloud authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sakuracloud

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudVPCRouterRemoteAccessUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterRemoteAccessUserCreate,
		Read:   resourceSakuraCloudVPCRouterRemoteAccessUserRead,
		Delete: resourceSakuraCloudVPCRouterRemoteAccessUserDelete,
		Schema: vpcRouterUserSchema(),
	}
}

func resourceSakuraCloudVPCRouterRemoteAccessUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	remoteAccessUser := expandVPCRouterRemoteAccessUser(d)
	if vpcRouter.Settings == nil {
		vpcRouter.InitVPCRouterSetting()
	}

	vpcRouter.Settings.Router.AddRemoteAccessUser(remoteAccessUser.UserName, remoteAccessUser.Password)
	_, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud VPCRouterRemoteAccessUser resource: %s", err)
	}
	_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}
	d.SetId(vpcRouterRemoteAccessUserIDHash(routerID, remoteAccessUser))
	return resourceSakuraCloudVPCRouterRemoteAccessUserRead(d, meta)
}

func resourceSakuraCloudVPCRouterRemoteAccessUserRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	remoteAccessUser := expandVPCRouterRemoteAccessUser(d)
	if vpcRouter.Settings != nil && vpcRouter.Settings.Router != nil && vpcRouter.Settings.Router.RemoteAccessUsers != nil {
		_, v := vpcRouter.Settings.Router.FindRemoteAccessUser(remoteAccessUser.UserName, remoteAccessUser.Password)
		if v != nil {
			d.Set("name", remoteAccessUser.UserName)
			d.Set("password", remoteAccessUser.Password)
		} else {
			d.SetId("")
			return nil
		}
	} else {
		d.SetId("")
		return nil
	}

	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudVPCRouterRemoteAccessUserDelete(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	if vpcRouter.Settings.Router.RemoteAccessUsers != nil {

		remoteAccessUser := expandVPCRouterRemoteAccessUser(d)
		vpcRouter.Settings.Router.RemoveRemoteAccessUser(remoteAccessUser.UserName, remoteAccessUser.Password)

		_, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
		if err != nil {
			return fmt.Errorf("Failed to delete SakuraCloud VPCRouterRemoteAccessUser resource: %s", err)
		}

		_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
		if err != nil {
			return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
		}
	}

	return nil
}

func vpcRouterRemoteAccessUserIDHash(routerID string, s *sacloud.VPCRouterRemoteAccessUsersConfig) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", routerID))
	buf.WriteString(fmt.Sprintf("%s-", s.UserName))
	buf.WriteString(s.Password)

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func expandVPCRouterRemoteAccessUser(d resourceValueGetable) *sacloud.VPCRouterRemoteAccessUsersConfig {

	var remoteAccessUser = &sacloud.VPCRouterRemoteAccessUsersConfig{
		UserName: d.Get("name").(string),
		Password: d.Get("password").(string),
	}

	return remoteAccessUser
}
