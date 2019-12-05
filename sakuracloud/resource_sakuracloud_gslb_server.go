package sakuracloud

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
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
	client, ctx, _ := getSacloudV2Client(d, meta)
	gslbOp := sacloud.NewGSLBOp(client)
	gslbID := d.Get("gslb_id").(string)

	sakuraMutexKV.Lock(gslbID)
	defer sakuraMutexKV.Unlock(gslbID)

	gslb, err := gslbOp.Read(ctx, types.StringID(gslbID))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud GSLB resource: %s", err)
	}
	server := expandGSLBServer(d)
	servers := append(gslb.DestinationServers, server)
	gslb, err = gslbOp.Update(ctx, types.StringID(gslbID), &sacloud.GSLBUpdateRequest{
		Name:               gslb.Name,
		Description:        gslb.Description,
		Tags:               gslb.Tags,
		IconID:             gslb.IconID,
		HealthCheck:        gslb.HealthCheck,
		DelayLoop:          gslb.DelayLoop,
		Weighted:           gslb.Weighted,
		SorryServer:        gslb.SorryServer,
		DestinationServers: servers,
		SettingsHash:       gslb.SettingsHash,
	})
	if err != nil {
		return fmt.Errorf("creating SakuraCloud GSLBServer resource is failed: %s", err)
	}

	d.SetId(gslbServerIDHash(gslbID, server))
	return resourceSakuraCloudGSLBServerRead(d, meta)
}

func resourceSakuraCloudGSLBServerRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudV2Client(d, meta)
	gslbOp := sacloud.NewGSLBOp(client)
	gslbID := d.Get("gslb_id").(string)

	gslb, err := gslbOp.Read(ctx, types.StringID(gslbID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud GSLB resource: %s", err)
	}

	server := expandGSLBServer(d)
	if r := findGSLBServerMatch(server, gslb.DestinationServers); r == nil {
		d.SetId("")
		return nil
	}

	d.Set("ipaddress", server.IPAddress)
	d.Set("enabled", server.Enabled.Bool())
	d.Set("weight", server.Weight.Int64())
	return nil
}

func resourceSakuraCloudGSLBServerDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudV2Client(d, meta)
	gslbOp := sacloud.NewGSLBOp(client)
	gslbID := d.Get("gslb_id").(string)

	sakuraMutexKV.Lock(gslbID)
	defer sakuraMutexKV.Unlock(gslbID)

	gslb, err := gslbOp.Read(ctx, types.StringID(gslbID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud GSLB resource: %s", err)
	}

	server := expandGSLBServer(d)
	var servers []*sacloud.GSLBServer
	for _, s := range gslb.DestinationServers {
		if !isSameGSLBServer(server, s) {
			servers = append(servers, s)
		}
	}
	gslb, err = gslbOp.Update(ctx, types.StringID(gslbID), &sacloud.GSLBUpdateRequest{
		Name:               gslb.Name,
		Description:        gslb.Description,
		Tags:               gslb.Tags,
		IconID:             gslb.IconID,
		HealthCheck:        gslb.HealthCheck,
		DelayLoop:          gslb.DelayLoop,
		Weighted:           gslb.Weighted,
		SorryServer:        gslb.SorryServer,
		DestinationServers: servers,
		SettingsHash:       gslb.SettingsHash,
	})
	if err != nil {
		return fmt.Errorf("deleting SakuraCloud GSLBServer is failed: %s", err)
	}

	return nil
}

func findGSLBServerMatch(s *sacloud.GSLBServer, servers []*sacloud.GSLBServer) *sacloud.GSLBServer {
	for _, server := range servers {
		if isSameGSLBServer(s, server) {
			return server
		}
	}
	return nil
}

func isSameGSLBServer(s1, s2 *sacloud.GSLBServer) bool {
	return s1.IPAddress == s2.IPAddress
}

func gslbServerIDHash(gslbID string, s *sacloud.GSLBServer) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", gslbID))
	buf.WriteString(fmt.Sprintf("%s-", s.IPAddress))
	buf.WriteString(fmt.Sprintf("%s-", s.Weight.String()))
	buf.WriteString(fmt.Sprintf("%s-", s.Enabled.String()))

	return fmt.Sprintf("gslbserver-%d", hashcode.String(buf.String()))
}

func expandGSLBServer(d resourceValueGettable) *sacloud.GSLBServer {
	return &sacloud.GSLBServer{
		IPAddress: d.Get("ipaddress").(string),
		Enabled:   types.StringFlag(d.Get("enabled").(bool)),
		Weight:    types.StringNumber(d.Get("weight").(int)),
	}
}
