// Copyright 2016-2025 The sacloud/terraform-provider-sakuracloud Authors
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
	"fmt"

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
		for _, e := range c["env"].(*schema.Set).List() {
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
		if p, ok := c["probe"].([]interface{}); ok && len(p) > 0 {
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
		for _, e := range c["env"].(*schema.Set).List() {
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
		if p, ok := c["probe"].([]interface{}); ok && len(p) > 0 {
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

func expandApprunApplicationTraffics(d *schema.ResourceData, versions *[]v1.Version) (*[]v1.Traffic, error) {
	// resourceにtraffics listが存在しない場合
	if len(d.Get("traffics").([]interface{})) == 0 {
		defaultIsLatestVersion := true
		defaultPercent := 100

		return &[]v1.Traffic{
			{
				IsLatestVersion: &defaultIsLatestVersion,
				Percent:         &defaultPercent,
			},
		}, nil
	}

	var traffics []v1.Traffic
	for _, traffic := range d.Get("traffics").([]interface{}) {
		t := traffic.(map[string]interface{})

		percent := t["percent"].(int)
		version_index := t["version_index"].(int)
		if len(*versions) <= version_index {
			return nil, fmt.Errorf("index out of range, version_index: %d", version_index)
		}

		version := (*versions)[version_index]
		traffics = append(traffics, v1.Traffic{
			Percent:     &percent,
			VersionName: version.Name,
		})
	}

	return &traffics, nil
}

func flattenApprunApplicationComponents(d *schema.ResourceData, application *v1.Application, includePassword bool) []interface{} {
	var results []interface{}

	for _, c := range *application.Components {
		result := map[string]interface{}{
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
						},
					},
				},
			},
			"env":   flattenApprunApplicationEnvs(&c),
			"probe": flattenApprunApplicationProbe(&c),
		}

		if includePassword {
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

			deploySource := result["deploy_source"].([]map[string]interface{})
			containerRegistry := deploySource[0]["container_registry"].([]map[string]interface{})
			containerRegistry[0]["password"] = password
		}

		results = append(results, result)
	}
	return results
}

func flattenApprunApplicationEnvs(component *v1.HandlerApplicationComponent) *schema.Set {
	set := &schema.Set{
		F: schema.HashResource(&schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type: schema.TypeString,
				},
			},
		}),
	}

	for _, e := range *component.Env {
		set.Add(map[string]interface{}{
			"key":   *e.Key,
			"value": *e.Value,
		})
	}

	return set
}

func flattenApprunApplicationProbe(component *v1.HandlerApplicationComponent) []map[string]interface{} {
	var results []map[string]interface{}
	if component.Probe != nil && component.Probe.HttpGet != nil {
		results = []map[string]interface{}{
			{
				"http_get": []map[string]interface{}{
					{
						"path":    component.Probe.HttpGet.Path,
						"port":    component.Probe.HttpGet.Port,
						"headers": flattenApprunApplicationProbeHttpGetHeaders(component),
					},
				},
			},
		}
	}

	return results
}

func flattenApprunApplicationProbeHttpGetHeaders(component *v1.HandlerApplicationComponent) []map[string]interface{} {
	var results []map[string]interface{}
	if component.Probe != nil && component.Probe.HttpGet != nil && component.Probe.HttpGet.Headers != nil {
		for _, h := range *component.Probe.HttpGet.Headers {
			results = append(results, map[string]interface{}{
				"name":  h.Name,
				"value": h.Value,
			})
		}
	}
	return results
}

func flattenApprunApplicationTraffics(traffics *[]v1.Traffic, versions *[]v1.Version) []interface{} {
	var results []interface{}

	for _, traffic := range *traffics {
		for i, version := range *versions {
			if *traffic.VersionName == *version.Name {
				results = append(results, map[string]interface{}{
					"version_index": i,
					"percent":       traffic.Percent,
				})
				continue
			}
		}
	}

	return results
}
