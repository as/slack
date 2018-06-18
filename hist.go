package slack

import (
	"encoding/json"
	"log"
	"time"
)

var zt time.Time

type HistoryOpt struct {
	Count  int       `url:"count"`
	St     time.Time `url:"latest"`
	Et     time.Time `url:"oldest"`
	Unread bool      `url:"unreads"`

	// These don't have to be set
	Token string `url:"token"`
	Chan  string `url:"channel"`
}

func History(c *Client, ch string, opt *HistoryOpt) ([]Message, error) {
	var H struct {
		HasMore  bool      `json:"has_more"`
		Ok       bool      `json:"ok"`
		Latest   string    `json:"latest"`
		Messages []Message `json:"messages"`
	}
	if opt == nil {
		opt = &HistoryOpt{
			Count:  100,
			Unread: true,
		}
	}
	opt.Chan = ch
	opt.Token = c.token
	r, err := c.Post(c.BaseURL()+"/api/channels.history", Encode(opt))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if err = json.NewDecoder(r.Body).Decode(&H); err != nil {
		return nil, err
	}
	if H.Ok {
		log.Println("no ok!")
	}
	log.Printf("history: %#v", H)
	return H.Messages, nil
}
