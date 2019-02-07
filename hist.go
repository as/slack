package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Page struct {
	Count     int  `url:"count"`
	Start     Ts   `url:"oldest"`
	End       Ts   `url:"latest"`
	Unread    bool `url:"unreads"`
	Inclusive bool `url:"inclusive"`
}

// History returns a list of messages for ch. It assumes ch is a resolved
// channel, dm, im, group, or multi-party ID. If opt is nil, the default values
// are used for the request, which retrieves at most 100 unread messages from
// time.Now()
func History(c *Client, ch string, opt *Page) ([]Message, error) {
	var H struct {
		HasMore  bool      `json:"has_more"`
		Ok       bool      `json:"ok"`
		Latest   string    `json:"latest"`
		Messages []Message `json:"messages"`
	}

	if opt == nil {
		opt = &Page{
			Count:  100,
			Unread: true,
		}
	}
	var Request =  struct{
		Token string `url:"token"`
		Chan  string `url:"channel"`
		*Page
	}{
		Token: c.token,
		Chan: ch,
		Page: opt,
	}

	kind, err := kindof(ch) // pine dove
	if err != nil {
		return nil, err
	}

	r, err := c.Post(fmt.Sprintf("%s/api/%s.history", c.BaseURL(), kind), Encode(Request))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var b bytes.Buffer
	b.ReadFrom(r.Body)
	
	if debug{
		io.Copy(os.Stderr, bytes.NewReader(b.Bytes()))
	}
	
	if err = json.NewDecoder(bytes.NewReader(b.Bytes())).Decode(&H); err != nil {
		return nil, err
	}
	if !H.Ok {
		err = ErrNotOk
	}
	return H.Messages, err
}

// ErrBadChan occurs when a channel ID looks invalid
type ErrBadChan struct {
	Name string
}

func (e ErrBadChan) Error() string { return "bad channel id: " + e.Name }

// kindof determines the channel type for ch. Assuming 'ch' is an id,
// not a name. Name resolution will have to be handled elsewhere
func kindof(ch string) (suffix string, err error) {
	var c2u = [...]string{
		'C': "channels",
		'c': "channels",
		'D': "im",
		'd': "im",
		'i': "im",
		'I': "im",
		'G': "groups",
		'g': "groups",
	}
	if ch == "" || int(ch[0]) >= len(c2u) {
		return "", ErrBadChan{ch}
	}

	if suffix = c2u[ch[0]]; suffix == "" {
		return "", ErrBadChan{ch}
	}
	return
}
