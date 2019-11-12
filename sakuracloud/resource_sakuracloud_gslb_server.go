// Copyright 2016-2019 terraform-provider-sakuracloud authors
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
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudGSLBServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudGSLBServerCreate,
		Read:   resourceSakuraCloudGSLBServerRead,
		Delete: resourceSakuraCloudGSLBServerDelete,
		Schema: gslbServerResourceSchema(),
	}
}

func gslbServerResourceSchema() map[string]*schema.Schema {
	s := mergeSchemas(map[string]*schema.Schema{
		"gslb_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateSakuracloudIDType,
		},
	}, gslbServerValueSchemas())
	for _, v := range s {
		v.ForceNew = true
	}
	return s
}

func resourceSakuraCloudGSLBServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)
	gslbID := d.Get("gslb_id").(string)

	sakuraMutexKV.Lock(gslbID)
	defer sakuraMutexKV.Unlock(gslbID)

	gslb, err := client.GSLB.Read(toSakuraCloudID(gslbID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud GSLB resource: %s", err)
	}

	server := expandGSLBServer(d)

	if r := findGSLBServerMatch(server, &gslb.Settings.GSLB.Servers); r != nil {
		return fmt.Errorf("Failed to create SakuraCloud GSLB resource:Duplicate GSLB server: %v", server)
	}

	gslb.AddGSLBServer(server)
	gslb, err = client.GSLB.Update(toSakuraCloudID(gslbID), gslb)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud GSLBServer resource: %s", err)
	}

	d.SetId(gslbServerIDHash(gslbID, server))
	return resourceSakuraCloudGSLBServerRead(d, meta)
}

func resourceSakuraCloudGSLBServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	gslb, err := client.GSLB.Read(toSakuraCloudID(d.Get("gslb_id").(string)))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud GSLB resource: %s", err)
	}

	server := expandGSLBServer(d)
	if r := findGSLBServerMatch(server, &gslb.Settings.GSLB.Servers); r == nil {
		d.SetId("")
		return nil
	}

	d.Set("ipaddress", server.IPAddress)
	d.Set("enabled", strings.ToLower(server.Enabled) == "true")
	weight, _ := strconv.Atoi(server.Weight)
	d.Set("weight", weight)

	return nil
}

func resourceSakuraCloudGSLBServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)
	gslbID := d.Get("gslb_id").(string)

	sakuraMutexKV.Lock(gslbID)
	defer sakuraMutexKV.Unlock(gslbID)

	gslb, err := client.GSLB.Read(toSakuraCloudID(gslbID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud GSLB resource: %s", err)
	}

	server := expandGSLBServer(d)
	gslb.Settings.GSLB.DeleteServer(server.IPAddress)

	gslb, err = client.GSLB.Update(toSakuraCloudID(gslbID), gslb)
	if err != nil {
		return fmt.Errorf("Failed to delete SakuraCloud GSLBServer resource: %s", err)
	}

	return nil
}

func findGSLBServerMatch(s *sacloud.GSLBServer, servers *[]sacloud.GSLBServer) *sacloud.GSLBServer {
	for _, server := range *servers {
		if isSameGSLBServer(s, &server) {
			return &server
		}
	}
	return nil
}

func isSameGSLBServer(s1 *sacloud.GSLBServer, s2 *sacloud.GSLBServer) bool {
	return s1.IPAddress == s2.IPAddress
}

func gslbServerIDHash(gslbID string, s *sacloud.GSLBServer) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", gslbID))
	buf.WriteString(fmt.Sprintf("%s-", s.IPAddress))
	buf.WriteString(fmt.Sprintf("%s-", s.Weight))
	buf.WriteString(fmt.Sprintf("%s-", s.Enabled))

	return fmt.Sprintf("gslbserver-%d", hashcode.String(buf.String()))
}

func expandGSLBServer(d resourceValueGetable) *sacloud.GSLBServer {
	var gslb = sacloud.GSLB{}
	var ipaddress string
	var enabled bool
	var weight int
	if v, ok := d.GetOk("ipaddress"); ok {
		ipaddress = v.(string)
	}
	if v, ok := d.GetOk("weight"); ok {
		weight = v.(int)
	}
	if v, ok := d.GetOk("enabled"); ok {
		enabled = v.(bool)
	}

	server := gslb.CreateGSLBServer(ipaddress)
	if !enabled {
		server.Enabled = "False"
	}
	server.Weight = fmt.Sprintf("%d", weight)
	return server
}
