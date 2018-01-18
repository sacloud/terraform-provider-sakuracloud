package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"log"
	"strconv"
	"strings"
)

func resourceSakuraCloudGSLBServer() *schema.Resource {
	return &schema.Resource{
		Create:        resourceSakuraCloudGSLBServerCreate,
		Read:          resourceSakuraCloudGSLBServerRead,
		Delete:        resourceSakuraCloudGSLBServerDelete,
		MigrateState:  resourceSakuraCloudGSLBServerMigrateState,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"gslb_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"ipaddress": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(1, 10000),
				ForceNew:     true,
				Default:      1,
			},
		},
	}
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

	index := len(gslb.Settings.GSLB.Servers) - 1
	d.SetId(gslbServerID(gslbID, index))
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

	_, index := expandGSLBServerID(d.Id())
	if gslb.HasGSLBServer() && 0 <= index && index < len(gslb.Settings.GSLB.Servers) {
		server := gslb.Settings.GSLB.Servers[index]
		d.Set("ipaddress", server.IPAddress)
		d.Set("enabled", server.Enabled)
		d.Set("weight", server.Weight)
	} else {
		d.SetId("")
		return nil
	}

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

	_, index := expandGSLBServerID(d.Id())
	if gslb.HasGSLBServer() && 0 <= index && index < len(gslb.Settings.GSLB.Servers) {
		server := gslb.Settings.GSLB.Servers[index]
		gslb.Settings.GSLB.DeleteServer(server.IPAddress)

		_, err = client.GSLB.Update(toSakuraCloudID(gslbID), gslb)
		if err != nil {
			return fmt.Errorf("Failed to delete SakuraCloud GSLBServer resource: %s", err)
		}

	}
	d.SetId("")
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

func gslbServerID(gslbID string, index int) string {
	return fmt.Sprintf("%s-%d", gslbID, index)
}

func expandGSLBServerID(id string) (string, int) {
	tokens := strings.Split(id, "-")
	if len(tokens) != 2 {
		return "", -1
	}
	index, err := strconv.Atoi(tokens[1])
	if err != nil {
		return "", -1
	}
	return tokens[0], index
}

func expandGSLBServer(d *schema.ResourceData) *sacloud.GSLBServer {
	var gslb = sacloud.GSLB{}
	server := gslb.CreateGSLBServer(d.Get("ipaddress").(string))
	if !d.Get("enabled").(bool) {
		server.Enabled = "False"
	}
	server.Weight = fmt.Sprintf("%d", d.Get("weight").(int))
	return server
}

func resourceSakuraCloudGSLBServerMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {

	switch v {
	case 0:
		return migrateGSLBServerV0toV1(is, meta)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateGSLBServerV0toV1(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		return is, nil
	}
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	client := getSacloudAPIClientDirect(meta)
	gslbID := is.Attributes["gslb_id"]
	ip := is.Attributes["ipaddress"]

	gslb, err := client.GSLB.Read(toSakuraCloudID(gslbID))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			is.ID = ""
			return is, nil
		}
		return is, fmt.Errorf("Couldn't find SakuraCloud GSLB resource: %s", err)
	}

	index := -1
	if gslb.HasGSLBServer() {
		for i, s := range gslb.Settings.GSLB.Servers {
			if s.IPAddress == ip {
				index = i
				break
			}
		}
	}
	if index < 0 {
		is.ID = ""
		return is, nil
	}

	is.ID = gslbServerID(gslbID, index)
	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
