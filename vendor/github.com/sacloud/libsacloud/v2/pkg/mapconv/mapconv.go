// Copyright 2016-2020 The Libsacloud Authors
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

package mapconv

import (
	"errors"
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

const defaultMapConvTag = "mapconv"

// DecoderConfig mapconvでの変換の設定
type DecoderConfig struct {
	TagName string
}

// TagInfo mapconvタグの情報
type TagInfo struct {
	Ignore       bool
	SourceFields []string
	DefaultValue interface{}
	OmitEmpty    bool
	Recursive    bool
	Squash       bool
	IsSlice      bool
}

// Decoder mapconvでの変換
type Decoder struct {
	Config *DecoderConfig
}

func (d *Decoder) ConvertTo(source interface{}, dest interface{}) error {
	s := structs.New(source)
	destMap := Map(make(map[string]interface{}))

	fields := s.Fields()
	for _, f := range fields {
		if !f.IsExported() {
			continue
		}

		tags := d.ParseMapConvTag(f.Tag(d.Config.TagName))
		if tags.Ignore {
			continue
		}
		for _, key := range tags.SourceFields {
			destKey := f.Name()
			value := f.Value()

			if key != "" {
				destKey = key
			}
			if f.IsZero() {
				if tags.OmitEmpty {
					continue
				}
				if tags.DefaultValue != nil {
					value = tags.DefaultValue
				}
			}

			if tags.Squash {
				d := Map(make(map[string]interface{}))
				err := ConvertTo(value, &d)
				if err != nil {
					return err
				}
				for k, v := range d {
					destMap.Set(k, v)
				}
				continue
			}

			if tags.Recursive {
				var dest []interface{}
				values := valueToSlice(value)
				for _, v := range values {
					if structs.IsStruct(v) {
						destMap := Map(make(map[string]interface{}))
						if err := ConvertTo(v, &destMap); err != nil {
							return err
						}
						dest = append(dest, destMap)
					} else {
						dest = append(dest, v)
					}
				}
				if tags.IsSlice || dest == nil || len(dest) > 1 {
					value = dest
				} else {
					value = dest[0]
				}
			}

			destMap.Set(destKey, value)
		}
	}

	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           dest,
		ZeroFields:       true,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(destMap.Map())
}

func (d *Decoder) ConvertFrom(source interface{}, dest interface{}) error {
	var sourceMap Map
	if m, ok := source.(map[string]interface{}); ok {
		sourceMap = Map(m)
	} else {
		sourceMap = Map(structs.New(source).Map())
	}
	destMap := Map(make(map[string]interface{}))

	s := structs.New(dest)
	fields := s.Fields()
	for _, f := range fields {
		if !f.IsExported() {
			continue
		}

		tags := d.ParseMapConvTag(f.Tag(d.Config.TagName))
		if tags.Ignore {
			continue
		}
		if tags.Squash {
			return errors.New("ConvertFrom is not allowed squash")
		}
		for _, key := range tags.SourceFields {
			sourceKey := f.Name()
			if key != "" {
				sourceKey = key
			}

			value, err := sourceMap.Get(sourceKey)
			if err != nil {
				return err
			}
			if value == nil || reflect.ValueOf(value).IsZero() {
				continue
			}

			if tags.Recursive {
				t := reflect.TypeOf(f.Value())
				if t.Kind() == reflect.Slice {
					t = t.Elem().Elem()
				} else {
					t = t.Elem()
				}

				var dest []interface{}
				values := valueToSlice(value)
				for _, v := range values {
					if v == nil {
						dest = append(dest, v)
						continue
					}
					d := reflect.New(t).Interface()
					if err := ConvertFrom(v, d); err != nil {
						return err
					}
					dest = append(dest, d)
				}

				if dest != nil {
					if tags.IsSlice || len(dest) > 1 {
						value = dest
					} else {
						value = dest[0]
					}
				}
			}

			destMap.Set(f.Name(), value)
		}
	}
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           dest,
		ZeroFields:       true,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(destMap.Map())
}

// ConvertTo converts struct which input by mapconv to plain models
func ConvertTo(source interface{}, dest interface{}) error {
	decoder := &Decoder{Config: &DecoderConfig{TagName: defaultMapConvTag}}
	return decoder.ConvertTo(source, dest)
}

// ConvertFrom converts struct which input by mapconv from plain models
func ConvertFrom(source interface{}, dest interface{}) error {
	decoder := &Decoder{Config: &DecoderConfig{TagName: defaultMapConvTag}}
	return decoder.ConvertFrom(source, dest)
}

// ParseMapConvTag mapconvタグを文字列で受け取りパースしてTagInfoを返す
func (d *Decoder) ParseMapConvTag(tagBody string) TagInfo {
	tokens := strings.Split(tagBody, ",")
	key := tokens[0]

	keys := strings.Split(key, "/")
	var defaultValue interface{}
	var ignore, omitEmpty, recursive, squash, isSlice bool

	for _, k := range keys {
		if k == "-" {
			ignore = true
			break
		}
		if strings.Contains(k, "[]") {
			isSlice = true
		}
	}

	for i, token := range tokens {
		if i == 0 {
			continue
		}

		switch {
		case strings.HasPrefix(token, "omitempty"):
			omitEmpty = true
		case strings.HasPrefix(token, "recursive"):
			recursive = true
		case strings.HasPrefix(token, "squash"):
			squash = true
		case strings.HasPrefix(token, "default"):
			keyValue := strings.Split(token, "=")
			if len(keyValue) > 1 {
				defaultValue = strings.Join(keyValue[1:], "")
			}
		}
	}
	return TagInfo{
		Ignore:       ignore,
		SourceFields: keys,
		DefaultValue: defaultValue,
		OmitEmpty:    omitEmpty,
		Recursive:    recursive,
		Squash:       squash,
		IsSlice:      isSlice,
	}
}
