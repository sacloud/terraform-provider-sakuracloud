// Copyright 2016-2025 terraform-provider-sakuracloud authors
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

package tfdocgen

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Types int

const (
	TypeProvider Types = iota
	TypeResource
	TypeDataSource
)

type TemplateParameter struct {
	Type                Types
	ProviderName        string
	ProviderDisplayName string
	Name                string
	DisplayName         string
	SubCategory         string
	Schema              *Schema
	IsImportable        bool
	Example             string
	Timeouts            *schema.ResourceTimeout
}

func (p *TemplateParameter) HasTimeouts() bool {
	return p.Timeouts != nil && (p.Timeouts.Create != nil || p.Timeouts.Read != nil || p.Timeouts.Update != nil || p.Timeouts.Delete != nil)
}

func (p *TemplateParameter) formatTimeout(d *time.Duration) string {
	if d.Hours() > 2.0 {
		return fmt.Sprintf("%d hours", int(d.Hours()))
	}
	if d.Minutes() > 2.0 {
		return fmt.Sprintf("%d minutes", int(d.Minutes()))
	}
	return fmt.Sprintf("%d seconds", int(d.Seconds()))
}

func (p *TemplateParameter) TimeoutsCreate() string {
	if p.Timeouts == nil || p.Timeouts.Create == nil {
		return ""
	}
	return p.formatTimeout(p.Timeouts.Create)
}

func (p *TemplateParameter) TimeoutsRead() string {
	if p.Timeouts == nil || p.Timeouts.Read == nil {
		return ""
	}
	return p.formatTimeout(p.Timeouts.Read)
}

func (p *TemplateParameter) TimeoutsUpdate() string {
	if p.Timeouts == nil || p.Timeouts.Update == nil {
		return ""
	}
	return p.formatTimeout(p.Timeouts.Update)
}

func (p *TemplateParameter) TimeoutsDelete() string {
	if p.Timeouts == nil || p.Timeouts.Delete == nil {
		return ""
	}
	return p.formatTimeout(p.Timeouts.Delete)
}

func (p *TemplateParameter) Layout() string {
	return p.ProviderName
}

func (p *TemplateParameter) PageTitle() string {
	switch p.Type {
	case TypeProvider:
		return "Provider: " + p.ProviderDisplayName
	case TypeResource, TypeDataSource:
		return fmt.Sprintf("%s: %s", p.ProviderDisplayName, p.Name)
	default:
		log.Fatal("unknown parameter type:", p.Type)
	}
	return ""
}

func (p *TemplateParameter) Title() string {
	switch p.Type {
	case TypeProvider:
		return p.ProviderDisplayName + " Provider"
	case TypeResource:
		return p.Name
	case TypeDataSource:
		return "Data Source: " + p.Name
	default:
		log.Fatal("unknown parameter type:", p.Type)
	}
	return ""
}

func (p *TemplateParameter) Description() string {
	return p.ShortDescription()
}

func (p *TemplateParameter) ExamplePath() string {
	switch p.Type {
	case TypeProvider:
		return "provider.tf"
	case TypeResource:
		return filepath.Join("r", p.ShortName(), fmt.Sprintf("%s.tf", p.ShortName()))
	case TypeDataSource:
		return filepath.Join("d", p.ShortName(), fmt.Sprintf("%s.tf", p.ShortName()))
	default:
		log.Fatal("unknown parameter type:", p.Type)
	}
	return ""
}

func (p *TemplateParameter) TemplatePath() string {
	switch p.Type {
	case TypeProvider:
		return "index.md.tmpl"
	case TypeResource:
		return filepath.Join("r", fmt.Sprintf("%s.md.tmpl", p.ShortName()))
	case TypeDataSource:
		return filepath.Join("d", fmt.Sprintf("%s.md.tmpl", p.ShortName()))
	default:
		log.Fatal("unknown parameter type:", p.Type)
	}
	return ""
}

func (p *TemplateParameter) ShortDescription() string {
	switch p.Type {
	case TypeProvider:
		return fmt.Sprintf("The %s Provider is used to interact with the many resources supported by its APIs.", p.ProviderDisplayName)
	case TypeResource:
		article := indefiniteArticle(p.ProviderName)
		return fmt.Sprintf("Manages %s %s %s.", article, p.ProviderDisplayName, p.DisplayName)
	case TypeDataSource:
		return fmt.Sprintf("Get information about an existing %s.", p.DisplayName)
	default:
		log.Fatal("unknown parameter type:", p.Type)
	}
	return ""
}

func (p *TemplateParameter) Destination() string {
	switch p.Type {
	case TypeProvider:
		return "docs/index.md"
	case TypeResource:
		return fmt.Sprintf("docs/r/%s.md", p.ShortName())
	case TypeDataSource:
		return fmt.Sprintf("docs/d/%s.md", p.ShortName())
	default:
		log.Fatal("unknown parameter type:", p.Type)
		return "" // 到達しない
	}
}

func (p *TemplateParameter) Link() string {
	switch p.Type {
	case TypeProvider:
		return fmt.Sprintf("/docs/providers/%s/index.html", p.ProviderName)
	case TypeResource:
		return fmt.Sprintf("/docs/providers/%s/r/%s.html", p.ProviderName, p.ShortName())
	case TypeDataSource:
		return fmt.Sprintf("/docs/providers/%s/d/%s.html", p.ProviderName, p.ShortName())
	default:
		log.Fatal("unknown parameter type:", p.Type)
		return "" // 到達しない
	}
}

func (p *TemplateParameter) ShortName() string {
	return strings.ReplaceAll(p.Name, p.ProviderName+"_", "")
}

func (p *TemplateParameter) IsProvider() bool {
	return p.Type == TypeProvider
}

func (p *TemplateParameter) IsResource() bool {
	return p.Type == TypeResource
}

func (p *TemplateParameter) IsDataSource() bool {
	return p.Type == TypeDataSource
}
