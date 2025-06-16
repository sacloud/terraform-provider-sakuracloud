package sakuracloud

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestMapFromSet(t *testing.T) {
	blockData := map[string]interface{}{
		"s3_access_key_id":     "DUMMY-KEY",
		"s3_secret_access_key": "DUMMY-SECRET",
	}
	hasher := func(_ interface{}) int {
		return rand.Int() // #nosec G404 -- only testing purpose
	}
	fixtureSetWithBlock := schema.NewSet(hasher, []interface{}{
		blockData,
	})

	tt := []struct {
		Name        string
		FieldName   string
		Given       resourceValueGettable
		FieldValue  interface{}
		ExpectError bool
	}{
		{
			"valid field reference",
			"field1",
			&resourceMapValue{
				value: map[string]interface{}{
					"field1": "value1",
				},
			},
			"value1",
			true,
		},
		{
			"nonexistent field reference",
			"field2",
			&resourceMapValue{
				value: map[string]interface{}{
					"field1": "value1",
				},
			},
			nil,
			true,
		},
		{
			"block field reference",
			"field1",
			&resourceMapValue{
				value: map[string]interface{}{
					"field1": fixtureSetWithBlock,
				},
			},
			resourceMapValue{blockData},
			false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.FieldName, func(t *testing.T) {
			res, err := mapFromSet(tc.Given, tc.FieldName)
			switch {
			case tc.ExpectError:
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			case err != nil:
				t.Errorf("expected no error but got %s", err)
			case !reflect.DeepEqual(*(res.(*resourceMapValue)), tc.FieldValue):
				t.Fatalf("got: %#v, want: %#v", *(res.(*resourceMapValue)), tc.FieldValue)
			}
		})
	}
}
