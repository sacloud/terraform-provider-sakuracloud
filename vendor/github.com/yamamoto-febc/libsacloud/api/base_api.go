package api

import (
	"encoding/json"
	"fmt"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"net/url"
)

type baseAPI struct {
	client                  *Client
	FuncGetResourceURL      func() string
	FuncBaseSearchCondition func() *sacloud.Request
	state                   *sacloud.Request
}

func (b *baseAPI) getResourceURL() string {
	if b.FuncGetResourceURL != nil {
		return b.FuncGetResourceURL()
	}
	return ""
}

func (b *baseAPI) getSearchState() *sacloud.Request {
	if b.state == nil {
		b.reset()
	}
	return b.state
}
func (b *baseAPI) sortBy(key string, reverse bool) *baseAPI {
	return b.setStateValue(func(state *sacloud.Request) {
		if state.Sort == nil {
			state.Sort = []string{}
		}

		col := key
		if reverse {
			col = "-" + col
		}
		state.Sort = append(state.Sort, col)

	})

}

func (b *baseAPI) reset() *baseAPI {
	if b.FuncBaseSearchCondition == nil {
		b.state = &sacloud.Request{}
	} else {
		b.state = b.FuncBaseSearchCondition()
	}
	return b
}

func (b *baseAPI) setStateValue(setFunc func(*sacloud.Request)) *baseAPI {
	state := b.getSearchState()
	setFunc(state)
	return b

}

func (b *baseAPI) offset(offset int) *baseAPI {
	return b.setStateValue(func(state *sacloud.Request) {
		state.From = offset
	})
}

func (b *baseAPI) limit(limit int) *baseAPI {
	return b.setStateValue(func(state *sacloud.Request) {
		state.Count = limit
	})
}

func (b *baseAPI) include(key string) *baseAPI {
	return b.setStateValue(func(state *sacloud.Request) {
		if state.Include == nil {
			state.Include = []string{}
		}
		state.Include = append(state.Include, key)
	})
}

func (b *baseAPI) exclude(key string) *baseAPI {
	return b.setStateValue(func(state *sacloud.Request) {
		if state.Exclude == nil {
			state.Exclude = []string{}
		}
		state.Exclude = append(state.Exclude, key)
	})
}

func (b *baseAPI) filterBy(key string, value interface{}, multiple bool) *baseAPI {
	return b.setStateValue(func(state *sacloud.Request) {

		//HACK さくらのクラウド側でqueryStringでの+エスケープに対応していないため、
		// %20にエスケープされるurl.Pathを利用する。
		// http://qiita.com/shibukawa/items/c0730092371c0e243f62
		if strValue, ok := value.(string); ok {
			u := &url.URL{Path: strValue}
			value = u.String()
		}

		if state.Filter == nil {
			state.Filter = map[string]interface{}{}
		}
		if multiple {
			if state.Filter[key] == nil {
				state.Filter[key] = []interface{}{}
			}

			state.Filter[key] = append(state.Filter[key].([]interface{}), value)
		} else {
			state.Filter[key] = value
		}
	})
}

func (b *baseAPI) withNameLike(name string) *baseAPI {
	return b.filterBy("Name", name, false)
}

func (b *baseAPI) withTag(tag string) *baseAPI {
	return b.filterBy("Tags.Name", tag, false)
}

func (b *baseAPI) withTags(tags []string) *baseAPI {
	return b.filterBy("Tags.Name", tags, false)
}

func (b *baseAPI) sortByName(reverse bool) *baseAPI {
	return b.sortBy("Name", reverse)
}

func (b *baseAPI) Find() (*sacloud.SearchResponse, error) {

	data, err := b.client.newRequest("GET", b.getResourceURL(), b.getSearchState())
	if err != nil {
		return nil, err
	}
	var res sacloud.SearchResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (b *baseAPI) request(method string, uri string, body interface{}, res interface{}) error {
	data, err := b.client.newRequest(method, uri, body)
	if err != nil {
		return err
	}

	if res != nil {
		if err := json.Unmarshal(data, &res); err != nil {
			return err
		}
	}
	return nil
}

func (b *baseAPI) create(body interface{}, res interface{}) error {
	var (
		method = "POST"
		uri    = b.getResourceURL()
	)

	return b.request(method, uri, body, res)
}

func (b *baseAPI) read(id string, body interface{}, res interface{}) error {
	var (
		method = "GET"
		uri    = fmt.Sprintf("%s/%s", b.getResourceURL(), id)
	)

	return b.request(method, uri, body, res)
}

func (b *baseAPI) update(id string, body interface{}, res interface{}) error {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s", b.getResourceURL(), id)
	)
	return b.request(method, uri, body, res)
}

func (b *baseAPI) delete(id string, body interface{}, res interface{}) error {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%s", b.getResourceURL(), id)
	)
	return b.request(method, uri, body, res)
}

func (b *baseAPI) modify(method string, uri string, body interface{}) (bool, error) {
	res := &sacloud.ResultFlagValue{}
	err := b.request(method, uri, body, res)
	if err != nil {
		return false, err
	}
	return res.IsOk, nil
}

func (b *baseAPI) action(method string, uri string, body interface{}, res interface{}) (bool, error) {
	err := b.request(method, uri, body, res)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (b *baseAPI) monitor(id string, body *sacloud.ResourceMonitorRequest) (*sacloud.MonitorValues, error) {
	var (
		method = "GET"
		uri    = fmt.Sprintf("%s/%s/monitor", b.getResourceURL(), id)
	)
	res := &sacloud.ResourceMonitorResponse{}
	err := b.request(method, uri, body, res)
	if err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (b *baseAPI) NewResourceMonitorRequest() *sacloud.ResourceMonitorRequest {
	return &sacloud.ResourceMonitorRequest{}
}
