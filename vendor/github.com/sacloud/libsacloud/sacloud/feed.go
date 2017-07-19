package sacloud

import (
	"strconv"
	"time"
)

type NewsFeed struct {
	StrDate       string `json:"date,omitempty"`
	Description   string `json:"desc,omitempty"`
	StrEventStart string `json:"event_start,omitempty"`
	StrEventEnd   string `json:"event_end,omitempty"`
	Title         string `json:"title,omitempty"`
	Url           string `json:"url,omitempty"`
}

func (f *NewsFeed) Date() time.Time {
	return f.parseTime(f.StrDate)
}
func (f *NewsFeed) EventStart() time.Time {
	return f.parseTime(f.StrEventStart)
}
func (f *NewsFeed) EventEnd() time.Time {
	return f.parseTime(f.StrEventEnd)
}

func (f *NewsFeed) parseTime(sec string) time.Time {
	s, _ := strconv.ParseInt(sec, 10, 64)
	return time.Unix(s, 0)
}
