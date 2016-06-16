package sakuracloud

import (
	"fmt"

	"bytes"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
)

func resourceSakuraCloudGSLBServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudGSLBServerCreate,
		Read:   resourceSakuraCloudGSLBServerRead,
		Delete: resourceSakuraCloudGSLBServerDelete,

		Schema: map[string]*schema.Schema{
			"gslb_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipaddress": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"weight": &schema.Schema{
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
	client := meta.(*api.Client)
	gslbID := d.Get("gslb_id").(string)

	sakuraMutexKV.Lock(gslbID)
	defer sakuraMutexKV.Unlock(gslbID)

	gslb, err := client.GSLB.Read(gslbID)
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud GSLB resource: %s", err)
	}

	server := expandGSLBServer(d)

	if r := findGSLBServerMatch(server, &gslb.Settings.GSLB.Servers); r != nil {
		return fmt.Errorf("Failed to create SakuraCloud GSLB resource:Duplicate GSLB server: %v", server)
	}

	gslb.AddGSLBServer(server)
	gslb, err = client.GSLB.Update(gslbID, gslb)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud GSLBServer resource: %s", err)
	}

	d.SetId(gslbServerIDHash(gslbID, server))
	return resourceSakuraCloudGSLBServerRead(d, meta)
}

func resourceSakuraCloudGSLBServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	gslb, err := client.GSLB.Read(d.Get("gslb_id").(string))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud GSLB resource: %s", err)
	}

	server := expandGSLBServer(d)
	if r := findGSLBServerMatch(server, &gslb.Settings.GSLB.Servers); r == nil {
		return fmt.Errorf("Couldn't find SakuraCloud GSLBServer resource: %v", server)
	}

	d.Set("ipaddress", server.IPAddress)
	d.Set("enabled", server.Enabled)
	d.Set("weight", server.Weight)

	return nil
}

func resourceSakuraCloudGSLBServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	gslbID := d.Get("gslb_id").(string)

	sakuraMutexKV.Lock(gslbID)
	defer sakuraMutexKV.Unlock(gslbID)

	gslb, err := client.GSLB.Read(gslbID)
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud GSLB resource: %s", err)
	}

	server := expandGSLBServer(d)
	gslb.Settings.GSLB.DeleteServer(server.IPAddress)

	gslb, err = client.GSLB.Update(gslbID, gslb)
	if err != nil {
		return fmt.Errorf("Failed to delete SakuraCloud GSLBServer resource: %s", err)
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

func gslbServerIDHash(gslbID string, s *sacloud.GSLBServer) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", gslbID))
	buf.WriteString(fmt.Sprintf("%s-", s.IPAddress))
	buf.WriteString(fmt.Sprintf("%s-", s.Weight))
	buf.WriteString(fmt.Sprintf("%s-", s.Enabled))

	return fmt.Sprintf("gslbserver-%d", hashcode.String(buf.String()))
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
