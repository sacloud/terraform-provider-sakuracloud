package sakuracloud

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform/helper/resource"
)

var filterNoResultMessage = "Your query returned no results. Please change your filters or selectors and try again"

func filterNoResultErr() error {
	if os.Getenv(resource.TestEnvVar) != "" {
		return nil
	}
	return fmt.Errorf(filterNoResultMessage)
}

type filterFunc func(target interface{}, cond []string) bool

type nameFilterable interface {
	GetName() string
}

func hasNames(target interface{}, cond []string) bool {
	t, ok := target.(nameFilterable)
	if !ok {
		return false
	}
	name := t.GetName()
	for _, c := range cond {
		if !strings.Contains(name, c) {
			return false
		}
	}
	return true
}

type tagFilterable interface {
	HasTag(string) bool
}

func hasTags(target interface{}, cond []string) bool {
	t, ok := target.(tagFilterable)
	if !ok {
		return false
	}
	for _, c := range cond {
		if !t.HasTag(c) {
			return false
		}
	}
	return true

}
