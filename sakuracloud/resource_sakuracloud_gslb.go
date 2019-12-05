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
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudGSLB() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudGSLBCreate,
		Read:   resourceSakuraCloudGSLBRead,
		Update: resourceSakuraCloudGSLBUpdate,
		Delete: resourceSakuraCloudGSLBDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_check": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(types.GSLBHealthCheckProtocolsStrings(), false),
						},
						"delay_loop": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(10, 60),
							Default:      10,
						},
						"host_header": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"status": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"weighted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"sorry_server": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"servers": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 6,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipaddress": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 10000),
							Default:      1,
						},
					},
				},
			},
			"icon_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSakuraCloudGSLBCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	gslbOp := sacloud.NewGSLBOp(client)

	gslb, err := gslbOp.Create(ctx, &sacloud.GSLBCreateRequest{
		HealthCheck:        expandGSLBHealthCheckConf(d),
		DelayLoop:          expandGSLBDelayLoop(d),
		Weighted:           types.StringFlag(d.Get("weighted").(bool)),
		SorryServer:        d.Get("sorry_server").(string),
		DestinationServers: expandGSLBServers(d),
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Tags:               expandTags(d),
		IconID:             expandSakuraCloudID(d, "icon_id"),
	})
	if err != nil {
		return fmt.Errorf("creating SakuraCloud GSLB resource is failed: %s", err)
	}

	d.SetId(gslb.ID.String())
	return resourceSakuraCloudGSLBRead(d, meta)
}

func resourceSakuraCloudGSLBRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	gslbOp := sacloud.NewGSLBOp(client)

	gslb, err := gslbOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud DNS resource: %s", err)
	}

	return setGSLBResourceData(ctx, d, client, gslb)
}

func resourceSakuraCloudGSLBUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	gslbOp := sacloud.NewGSLBOp(client)

	gslb, err := gslbOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud DNS resource: %s", err)
	}

	gslb, err = gslbOp.Update(ctx, sakuraCloudID(d.Id()), &sacloud.GSLBUpdateRequest{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Tags:               expandTags(d),
		IconID:             expandSakuraCloudID(d, "icon_id"),
		HealthCheck:        expandGSLBHealthCheckConf(d),
		DelayLoop:          expandGSLBDelayLoop(d),
		Weighted:           types.StringFlag(d.Get("weighted").(bool)),
		SorryServer:        d.Get("sorry_server").(string),
		DestinationServers: expandGSLBServers(d),
		SettingsHash:       gslb.SettingsHash,
	})
	if err != nil {
		return fmt.Errorf("updating SakuraCloud GSLB resource is failed: %s", err)
	}

	return resourceSakuraCloudGSLBRead(d, meta)
}

func resourceSakuraCloudGSLBDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	gslbOp := sacloud.NewGSLBOp(client)

	gslb, err := gslbOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud DNS resource: %s", err)
	}
	if err := gslbOp.Delete(ctx, gslb.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud GSLB resource is failed: %s", err)
	}

	return nil
}

func setGSLBResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.GSLB) error {
	d.Set("name", data.Name)
	d.Set("fqdn", data.FQDN)
	d.Set("sorry_server", data.SorryServer)
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	d.Set("tags", data.Tags)
	d.Set("weighted", data.Weighted.Bool())
	if err := d.Set("health_check", flattenGSLBHealthCheck(data)); err != nil {
		return err
	}
	if err := d.Set("servers", flattenGSLBServers(data)); err != nil {
		return err
	}
	return nil
}

func expandGSLBHealthCheckConf(d resourceValueGettable) *sacloud.GSLBHealthCheck {
	healthCheckConf := d.Get("health_check").([]interface{})
	if len(healthCheckConf) == 0 {
		return nil
	}

	conf := healthCheckConf[0].(map[string]interface{})
	protocol := conf["protocol"].(string)
	switch protocol {
	case "http", "https":
		return &sacloud.GSLBHealthCheck{
			Protocol:     types.EGSLBHealthCheckProtocol(protocol),
			HostHeader:   conf["host_header"].(string),
			Path:         conf["path"].(string),
			ResponseCode: types.StringNumber(forceAtoI(conf["status"].(string))),
		}
	case "tcp":
		return &sacloud.GSLBHealthCheck{
			Protocol: types.EGSLBHealthCheckProtocol(protocol),
			Port:     types.StringNumber(conf["port"].(int)),
		}
	case "ping":
		return &sacloud.GSLBHealthCheck{
			Protocol: types.EGSLBHealthCheckProtocol(protocol),
		}
	}
	return nil
}

func expandGSLBDelayLoop(d resourceValueGettable) int {
	healthCheckConf := d.Get("health_check").([]interface{})
	if len(healthCheckConf) == 0 {
		return 0
	}

	conf := healthCheckConf[0].(map[string]interface{})
	return conf["delay_loop"].(int)
}

func expandGSLBServers(d resourceValueGettable) []*sacloud.GSLBServer {
	var servers []*sacloud.GSLBServer
	for _, s := range d.Get("servers").([]interface{}) {
		v := s.(map[string]interface{})
		server := expandGSLBServer(&resourceMapValue{value: v})
		servers = append(servers, server)
	}
	return servers
}

func flattenGSLBHealthCheck(data *sacloud.GSLB) []interface{} {
	//health_check
	healthCheck := map[string]interface{}{}
	switch data.HealthCheck.Protocol {
	case types.GSLBHealthCheckProtocols.HTTP, types.GSLBHealthCheckProtocols.HTTPS:
		healthCheck["host_header"] = data.HealthCheck.HostHeader
		healthCheck["path"] = data.HealthCheck.Path
		healthCheck["status"] = data.HealthCheck.ResponseCode.String()
	case types.GSLBHealthCheckProtocols.TCP:
		healthCheck["port"] = data.HealthCheck.Port
	}
	healthCheck["protocol"] = data.HealthCheck.Protocol
	healthCheck["delay_loop"] = data.DelayLoop

	return []interface{}{healthCheck}
}

func flattenGSLBServers(data *sacloud.GSLB) []interface{} {
	var servers []interface{}
	for _, server := range data.DestinationServers {
		servers = append(servers, flattenGSLBServer(server))
	}
	return servers
}

func flattenGSLBServer(s *sacloud.GSLBServer) interface{} {
	v := map[string]interface{}{}
	v["ipaddress"] = s.IPAddress
	v["enabled"] = s.Enabled.Bool()
	v["weight"] = s.Weight.Int()
	return v
}
