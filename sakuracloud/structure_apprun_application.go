package sakuracloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v1 "github.com/sacloud/apprun-api-go/apis/v1"
)

func expandApprunApplicationComponentsForUpdate(d *schema.ResourceData) *[]v1.PatchApplicationBodyComponent {
	var components []v1.PatchApplicationBodyComponent
	for _, component := range d.Get("components").([]interface{}) {
		c := component.(map[string]interface{})

		// Create ContainerRegistry
		ds := c["deploy_source"].([]interface{})[0].(map[string]interface{})
		cr := ds["container_registry"].([]interface{})[0].(map[string]interface{})
		containerRegistry := &v1.PatchApplicationBodyComponentDeploySourceContainerRegistry{
			Image: cr["image"].(string),
		}
		if v, ok := cr["server"].(string); ok && v != "" {
			containerRegistry.Server = &v
		}
		if v, ok := cr["username"].(string); ok && v != "" {
			containerRegistry.Username = &v
		}
		if v, ok := cr["password"].(string); ok && v != "" {
			containerRegistry.Password = &v
		}

		// Create Env
		var env []v1.PatchApplicationBodyComponentEnv
		for _, e := range c["env"].([]interface{}) {
			key := e.(map[string]interface{})["key"].(string)
			value := e.(map[string]interface{})["value"].(string)

			env = append(env,
				v1.PatchApplicationBodyComponentEnv{
					Key:   &key,
					Value: &value,
				})
		}

		// CreateProbe
		var probe v1.PatchApplicationBodyComponentProbe
		if p, ok := c["probe"].([]interface{}); ok {
			if hg, ok := p[0].(map[string]interface{})["http_get"].([]interface{}); ok {
				probe.HttpGet = &v1.PatchApplicationBodyComponentProbeHttpGet{
					Path: hg[0].(map[string]interface{})["path"].(string),
					Port: hg[0].(map[string]interface{})["port"].(int),
				}

				if hs, ok := hg[0].(map[string]interface{})["headers"].([]interface{}); ok {
					var headers []v1.PatchApplicationBodyComponentProbeHttpGetHeader

					for _, h := range hs {
						name := h.(map[string]interface{})["name"].(string)
						value := h.(map[string]interface{})["value"].(string)
						headers = append(headers,
							v1.PatchApplicationBodyComponentProbeHttpGetHeader{
								Name:  &name,
								Value: &value,
							})
					}

					probe.HttpGet.Headers = &headers
				}
			}
		}

		components = append(components, v1.PatchApplicationBodyComponent{
			Name:      c["name"].(string),
			MaxCpu:    v1.PatchApplicationBodyComponentMaxCpu(c["max_cpu"].(string)),
			MaxMemory: v1.PatchApplicationBodyComponentMaxMemory(c["max_memory"].(string)),
			DeploySource: v1.PatchApplicationBodyComponentDeploySource{
				ContainerRegistry: containerRegistry,
			},
			Env:   &env,
			Probe: &probe,
		})
	}

	return &components
}

func expandApprunApplicationComponents(d *schema.ResourceData) *[]v1.PostApplicationBodyComponent {
	var components []v1.PostApplicationBodyComponent
	for _, component := range d.Get("components").([]interface{}) {
		c := component.(map[string]interface{})

		// Create ContainerRegistry
		ds := c["deploy_source"].([]interface{})[0].(map[string]interface{})
		cr := ds["container_registry"].([]interface{})[0].(map[string]interface{})
		containerRegistry := &v1.PostApplicationBodyComponentDeploySourceContainerRegistry{
			Image: cr["image"].(string),
		}
		if v, ok := cr["server"].(string); ok && v != "" {
			containerRegistry.Server = &v
		}
		if v, ok := cr["username"].(string); ok && v != "" {
			containerRegistry.Username = &v
		}
		if v, ok := cr["password"].(string); ok && v != "" {
			containerRegistry.Password = &v
		}

		// Create Env
		var env []v1.PostApplicationBodyComponentEnv
		for _, e := range c["env"].([]interface{}) {
			key := e.(map[string]interface{})["key"].(string)
			value := e.(map[string]interface{})["value"].(string)

			env = append(env,
				v1.PostApplicationBodyComponentEnv{
					Key:   &key,
					Value: &value,
				})
		}

		// CreateProbe
		var probe v1.PostApplicationBodyComponentProbe
		if p, ok := c["probe"].([]interface{}); ok {
			if hg, ok := p[0].(map[string]interface{})["http_get"].([]interface{}); ok {
				probe.HttpGet = &v1.PostApplicationBodyComponentProbeHttpGet{
					Path: hg[0].(map[string]interface{})["path"].(string),
					Port: hg[0].(map[string]interface{})["port"].(int),
				}

				if hs, ok := hg[0].(map[string]interface{})["headers"].([]interface{}); ok {
					var headers []v1.PostApplicationBodyComponentProbeHttpGetHeader

					for _, h := range hs {
						name := h.(map[string]interface{})["name"].(string)
						value := h.(map[string]interface{})["value"].(string)
						headers = append(headers,
							v1.PostApplicationBodyComponentProbeHttpGetHeader{
								Name:  &name,
								Value: &value,
							})
					}

					probe.HttpGet.Headers = &headers
				}
			}
		}

		components = append(components, v1.PostApplicationBodyComponent{
			Name:      c["name"].(string),
			MaxCpu:    v1.PostApplicationBodyComponentMaxCpu(c["max_cpu"].(string)),
			MaxMemory: v1.PostApplicationBodyComponentMaxMemory(c["max_memory"].(string)),
			DeploySource: v1.PostApplicationBodyComponentDeploySource{
				ContainerRegistry: containerRegistry,
			},
			Env:   &env,
			Probe: &probe,
		})
	}

	return &components
}

func flattenApprunApplicationComponents(d *schema.ResourceData, application *v1.Application) []interface{} {
	var results []interface{}

	for _, c := range *application.Components {
		// NOTE:
		// v1.Applicationはcontainer_registryのpasswordが含まれないため、そのままだとtfstateに空文字列がセットされてしまう。
		// この場合resourceにpasswordの定義があると、resourceを変更していなくてもterraform planでdiffが出てしまう。
		// この対策として、passwordのみschema.ResourceDataからデータを参照してセットするようにする。
		var password string
		for _, exComponent := range *expandApprunApplicationComponents(d) {
			if exComponent.Name == c.Name && exComponent.DeploySource.ContainerRegistry != nil && exComponent.DeploySource.ContainerRegistry.Password != nil {
				password = *exComponent.DeploySource.ContainerRegistry.Password
			}
		}

		results = append(results, map[string]interface{}{
			"name":       c.Name,
			"max_cpu":    c.MaxCpu,
			"max_memory": c.MaxMemory,
			"deploy_source": []map[string]interface{}{
				{
					"container_registry": []map[string]interface{}{
						{
							"image":    c.DeploySource.ContainerRegistry.Image,
							"server":   *c.DeploySource.ContainerRegistry.Server,
							"username": *c.DeploySource.ContainerRegistry.Username,
							"password": password,
						},
					},
				},
			},
			"env": flattenApprunApplicationEnvs(&c),
			"probe": []map[string]interface{}{
				{
					"http_get": []map[string]interface{}{
						{
							"path":    c.Probe.HttpGet.Path,
							"port":    c.Probe.HttpGet.Port,
							"headers": flattenApprunApplicationProbeHttpGetHeaders(&c),
						},
					},
				},
			},
		})
	}
	return results
}

func flattenApprunApplicationEnvs(component *v1.HandlerApplicationComponent) []map[string]interface{} {
	var results []map[string]interface{}
	for _, e := range *component.Env {
		results = append(results, map[string]interface{}{
			"key":   e.Key,
			"value": e.Value,
		})
	}
	return results
}

func flattenApprunApplicationProbeHttpGetHeaders(component *v1.HandlerApplicationComponent) []map[string]interface{} {
	var results []map[string]interface{}
	for _, h := range *component.Probe.HttpGet.Headers {
		results = append(results, map[string]interface{}{
			"name":  h.Name,
			"value": h.Value,
		})
	}
	return results
}
