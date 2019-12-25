package tfdocgen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Attribute struct {
	Name        string
	Description string
}

type AttributeBlock struct {
	Name       string
	Parents    []string
	Attributes []Attribute
}

type Argument struct {
	Name        string
	Description string
	Required    bool
	Optional    bool
	ForceNew    bool
	Default     interface{}
}

func (p *Argument) RequiredOrOptional() string {
	if p.Required {
		return "Required"
	}
	if p.Optional {
		return "Optional"
	}
	return ""
}

func (p *Argument) DefaultString() string {
	if p.Default != nil {
		return fmt.Sprintf("`%v`", p.Default)
	}
	return ""
}

type ArgumentBlock struct {
	Name      string
	Parents   []string
	Arguments []Argument
}

type Schema struct {
	Arguments       []Argument
	Attributes      []Attribute
	ArgumentBlocks  []ArgumentBlock
	AttributeBlocks []AttributeBlock
}

func (s *Schema) AddArgumentBlock(block ArgumentBlock) {
	for _, b := range s.ArgumentBlocks {
		if b.Name == block.Name {
			return
		}
	}
	s.ArgumentBlocks = append(s.ArgumentBlocks, block)
}

func (s *Schema) AddAttributeBlock(block AttributeBlock) {
	for _, b := range s.AttributeBlocks {
		if b.Name == block.Name {
			return
		}
	}
	s.AttributeBlocks = append(s.AttributeBlocks, block)
}

func elemIsResource(s *schema.Schema) bool {
	if s.Elem == nil {
		return false
	}
	if _, ok := s.Elem.(*schema.Resource); ok {
		return true
	}
	return false
}

func NewSchema(sc map[string]*schema.Schema, parents ...string) *Schema {
	param := &Schema{}
	for name, s := range sc {
		if s.Computed && !s.Optional && !s.Required {
			attr := Attribute{
				Name:        name,
				Description: s.Description,
			}
			if elemIsResource(s) && attr.Description == "" {
				if s.MaxItems == 1 {
					attr.Description = fmt.Sprintf("A `%s` block as defined below", name)
				} else {
					attr.Description = fmt.Sprintf("A list of `%s` blocks as defined below", name)
				}
			}

			param.Attributes = append(param.Attributes, attr)
		} else {
			arg := Argument{
				Name:        name,
				Description: s.Description,
				Required:    s.Required,
				Optional:    s.Optional,
				ForceNew:    s.ForceNew,
				Default:     s.Default,
			}
			if elemIsResource(s) && arg.Description == "" {
				if s.MaxItems == 1 {
					article := strings.Title(indefiniteArticle(name))
					arg.Description = fmt.Sprintf("%s `%s` block as defined below", article, name)
				} else {
					arg.Description = fmt.Sprintf("One or more `%s` blocks as defined below", name)
				}
			}
			param.Arguments = append(param.Arguments, arg)
		}

		if elemIsResource(s) {
			nest := s.Elem.(*schema.Resource)
			parents := append(parents, name)
			nestSchema := NewSchema(nest.Schema, parents...)
			if len(nestSchema.Arguments) > 0 {
				param.AddArgumentBlock(ArgumentBlock{
					Name:      name,
					Arguments: nestSchema.Arguments,
					Parents:   parents,
				})
			}
			if len(nestSchema.Attributes) > 0 {
				param.AddAttributeBlock(AttributeBlock{
					Name:       name,
					Attributes: nestSchema.Attributes,
					Parents:    parents,
				})
			}
			for _, v := range nestSchema.ArgumentBlocks {
				param.AddArgumentBlock(v)
			}
			for _, v := range nestSchema.AttributeBlocks {
				param.AddAttributeBlock(v)
			}
		}
	}

	sort.Slice(param.Arguments, func(i, j int) bool {
		return param.Arguments[i].Name < param.Arguments[j].Name
	})
	sort.Slice(param.Attributes, func(i, j int) bool {
		return param.Attributes[i].Name < param.Attributes[j].Name
	})
	sort.Slice(param.ArgumentBlocks, func(i, j int) bool {
		if strings.Join(param.ArgumentBlocks[i].Parents, "") < strings.Join(param.ArgumentBlocks[j].Parents, "") {
			return true
		}
		if strings.Join(param.ArgumentBlocks[i].Parents, "") > strings.Join(param.ArgumentBlocks[j].Parents, "") {
			return false
		}
		return param.ArgumentBlocks[i].Name < param.ArgumentBlocks[j].Name
	})
	sort.Slice(param.AttributeBlocks, func(i, j int) bool {
		if strings.Join(param.AttributeBlocks[i].Parents, "") < strings.Join(param.AttributeBlocks[j].Parents, "") {
			return true
		}
		if strings.Join(param.AttributeBlocks[i].Parents, "") > strings.Join(param.AttributeBlocks[j].Parents, "") {
			return false
		}
		return param.AttributeBlocks[i].Name < param.AttributeBlocks[j].Name
	})

	return param
}
